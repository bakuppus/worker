package listeners

import (
	"context"
	"fmt"
	permcache "github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/common/sentry"
	"github.com/TicketsBot/worker"
	"github.com/TicketsBot/worker/bot/command"
	context2 "github.com/TicketsBot/worker/bot/command/context"
	"github.com/TicketsBot/worker/bot/command/manager"
	"github.com/TicketsBot/worker/bot/command/registry"
	"github.com/TicketsBot/worker/bot/customisation"
	"github.com/TicketsBot/worker/bot/metrics/prometheus"
	"github.com/TicketsBot/worker/bot/metrics/statsd"
	"github.com/TicketsBot/worker/bot/utils"
	"github.com/TicketsBot/worker/i18n"
	"github.com/rxdn/gdl/gateway/payloads/events"
	"github.com/rxdn/gdl/objects/interaction"
	"golang.org/x/sync/errgroup"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	channelPattern = regexp.MustCompile(`<#(\d+)>`)
	userPattern    = regexp.MustCompile(`<@!?(\d+)>`)
	rolePattern    = regexp.MustCompile(`<@&(\d+)>`)
)

func GetCommandListener() func(*worker.Context, *events.MessageCreate) {
	commandManager := new(manager.CommandManager)
	commandManager.RegisterCommands()
	commandManager.RunSetupFuncs()

	return func(worker *worker.Context, e *events.MessageCreate) {
		if e.Author.Bot {
			return
		}

		// Ignore commands in DMs
		if e.GuildId == 0 {
			return
		}

		e.Member.User = e.Author

		// fmt.Sprintf is twice as slow!
		mentionPrefix := "<@" + strconv.FormatUint(worker.BotId, 10) + ">"

		var usedPrefix string
		if strings.HasPrefix(e.Content, mentionPrefix) {
			usedPrefix = mentionPrefix
		} else if strings.HasPrefix(strings.ToLower(e.Content), utils.DefaultPrefix) {
			usedPrefix = utils.DefaultPrefix
		} else {
			return
		}

		content := strings.TrimPrefix(e.Content, usedPrefix)
		content = strings.TrimSpace(content)

		split := strings.Split(content, " ")
		root := split[0]

		args := make([]string, 0)
		if len(split) > 1 {
			for _, arg := range split[1:] {
				if arg != "" {
					args = append(args, arg)
				}
			}
		}

		var c, rootCmd registry.Command
		for _, cmd := range commandManager.GetCommands() {
			if strings.ToLower(cmd.Properties().Name) == strings.ToLower(root) || contains(cmd.Properties().Aliases, strings.ToLower(root)) {
				parent := cmd
				rootCmd = cmd
				index := 0

				for {
					if len(args) > index {
						childName := args[index]
						found := false

						for _, child := range parent.Properties().Children {
							if strings.ToLower(child.Properties().Name) == strings.ToLower(childName) || contains(child.Properties().Aliases, strings.ToLower(childName)) {
								parent = child
								found = true
								index++
							}
						}

						if !found {
							break
						}
					} else {
						break
					}
				}

				var childArgs []string
				if len(args) > 0 {
					childArgs = args[index:]
				}

				args = childArgs
				c = parent
			}
		}

		if c == nil {
			return
		}

		properties := c.Properties()
		if properties.MainBotOnly && worker.IsWhitelabel {
			return
		}

		userPermissionLevel, err := permcache.GetPermissionLevel(utils.ToRetriever(worker), e.Member, e.GuildId)
		if err != nil {
			sentry.Error(err)
			return
		}

		var blacklisted bool
		var premiumTier premium.PremiumTier

		group, _ := errgroup.WithContext(context.Background())

		// get blacklisted
		group.Go(func() (err error) {
			// If e.Member is zero, it does not matter, as it is not checked if the command is not executed in a guild
			blacklisted, err = utils.IsBlacklisted(e.GuildId, e.Author.Id, e.Member, userPermissionLevel)
			return
		})

		// get premium tier
		group.Go(func() (err error) {
			premiumTier, err = utils.PremiumClient.GetTierByGuildId(e.GuildId, true, worker.Token, worker.RateLimiter)
			return
		})

		if err := group.Wait(); err != nil {
			sentry.Error(err)
			return
		}

		// Ensure user isn't blacklisted
		if blacklisted {
			utils.ReactWithCross(worker, e.ChannelId, e.Id)
			return
		}

		ctx := context2.NewMessageContext(worker, e.Message, args, premiumTier, userPermissionLevel)

		// Redirect user to slash commands
		if !properties.AdminOnly && !properties.HelperOnly {
			commands, err := command.LoadCommandIds(worker, worker.BotId)
			if err != nil {
				sentry.Error(err)
				return
			}

			commandName := strings.ToLower(rootCmd.Properties().Name)

			commandMention := "`COMMAND NOT FOUND`"
			if id, ok := commands[commandName]; ok {
				commandMention = fmt.Sprintf(`</%s:%d>`, commandName, id)
			}

			ctx.Reject()
			ctx.Reply(customisation.Red, i18n.Error, i18n.MessageInteractionSwitch, commandMention)
			return
		}

		if properties.InteractionOnly || rootCmd.Properties().InteractionOnly {
			ctx.Reject()
			ctx.Reply(customisation.Red, i18n.Error, i18n.MessageInteractionOnly, rootCmd.Properties().Name)
			return
		}

		if properties.PermissionLevel > userPermissionLevel {
			ctx.Reject()
			ctx.Reply(customisation.Red, i18n.Error, i18n.MessageNoPermission)
			return
		}

		if properties.AdminOnly && !utils.IsBotAdmin(e.Author.Id) {
			ctx.Reject()
			ctx.Reply(customisation.Red, i18n.Error, i18n.MessageOwnerOnly)
			return
		}

		if properties.HelperOnly && !utils.IsBotHelper(e.Author.Id) {
			ctx.Reject()
			ctx.Reply(customisation.Red, i18n.Error, i18n.MessageNoPermission)
			return
		}

		if properties.PremiumOnly && premiumTier == premium.None {
			ctx.Reject()
			ctx.Reply(customisation.Red, i18n.TitlePremiumOnly, i18n.MessagePremium)
			return
		}

		// parse args
		parsedArguments := make([]interface{}, len(properties.Arguments))

		var argsIndex int
		for i, argument := range properties.Arguments {
			if !argument.MessageCompatible {
				parsedArguments[i] = nil
			}

			if argsIndex >= len(args) {
				if argument.Required {
					ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
					return
				}

				continue
			}

			// TODO: translate messages
			switch argument.Type {
			case interaction.OptionTypeString:
				parsedArguments[i] = strings.Join(args[argsIndex:], " ")
				argsIndex = len(args)
			case interaction.OptionTypeInteger:
				//goland:noinspection GoNilness
				raw := args[argsIndex]
				value, err := strconv.Atoi(raw)
				if err != nil {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = (*int)(nil)
						continue
					}
				}

				parsedArguments[i] = value
				argsIndex++
			case interaction.OptionTypeBoolean:
				//goland:noinspection GoNilness
				raw := args[argsIndex]
				value, err := strconv.ParseBool(raw)
				if err != nil {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = (*bool)(nil)
						continue
					}
				}

				parsedArguments[i] = value
				argsIndex++
			case interaction.OptionTypeUser:
				//goland:noinspection GoNilness
				match := userPattern.FindStringSubmatch(args[argsIndex])
				if len(match) < 2 {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = nil
						continue
					}
				}

				if userId, err := strconv.ParseUint(match[1], 10, 64); err == nil {
					parsedArguments[i] = userId
					argsIndex++
				} else {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = nil
						continue
					}
				}
			case interaction.OptionTypeChannel:
				//goland:noinspection GoNilness
				match := channelPattern.FindStringSubmatch(args[argsIndex])
				if len(match) < 2 {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = nil
						continue
					}
				}

				if channelId, err := strconv.ParseUint(match[1], 10, 64); err == nil {
					parsedArguments[i] = channelId
					argsIndex++
				} else {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = nil
						continue
					}
				}
			case interaction.OptionTypeRole:
				//goland:noinspection GoNilness
				match := rolePattern.FindStringSubmatch(args[argsIndex])
				if len(match) < 2 {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = nil
						continue
					}
				}

				if roleId, err := strconv.ParseUint(match[1], 10, 64); err == nil {
					parsedArguments[i] = roleId
					argsIndex++
				} else {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = nil
						continue
					}
				}
			case interaction.OptionTypeMentionable:
				var snowflake string

				// First, check for role
				//goland:noinspection GoNilness
				roleMatch := rolePattern.FindStringSubmatch(args[argsIndex])
				if len(roleMatch) < 2 {
					// Then, check for user
					userMatch := userPattern.FindStringSubmatch(args[argsIndex])
					if len(userMatch) < 2 {
						if argument.Required {
							ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
							return
						} else {
							parsedArguments[i] = nil
							continue
						}
					} else {
						snowflake = userMatch[1]
					}
				} else {
					snowflake = roleMatch[1]
				}

				if id, err := strconv.ParseUint(snowflake, 10, 64); err == nil {
					parsedArguments[i] = id
					argsIndex++
				} else {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = nil
						continue
					}
				}
			case interaction.OptionTypeNumber:
				//goland:noinspection GoNilness
				raw := args[argsIndex]
				value, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					if argument.Required {
						ctx.Reply(customisation.Red, i18n.Error, argument.InvalidMessage)
						return
					} else {
						parsedArguments[i] = (*float64)(nil)
						continue
					}
				}

				parsedArguments[i] = value
				argsIndex++
			default:
				ctx.HandleError(fmt.Errorf("unknown argument type: %d", argument.Type))
			}
		}

		e.Member.User = e.Author

		valueArgs := make([]reflect.Value, len(parsedArguments)+1)
		valueArgs[0] = reflect.ValueOf(&ctx)

		fn := reflect.TypeOf(c.GetExecutor())
		for i, arg := range parsedArguments {
			var value reflect.Value
			if properties.Arguments[i].Required && arg != nil {
				value = reflect.ValueOf(arg)
			} else {
				if arg == nil {
					value = reflect.ValueOf(arg)
				} else {
					value = reflect.New(reflect.TypeOf(arg))
					tmp := value.Elem()
					tmp.Set(reflect.ValueOf(arg))
				}
			}

			if !value.IsValid() {
				value = reflect.New(fn.In(i + 1)).Elem()
			}

			valueArgs[i+1] = value
		}

		go reflect.ValueOf(c.GetExecutor()).Call(valueArgs)
		statsd.Client.IncrementKey(statsd.KeyCommands)
		prometheus.LogCommand(e.GuildId, rootCmd.Properties().Name)

		utils.DeleteAfter(worker, e.ChannelId, e.Id, utils.DeleteAfterSeconds)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
