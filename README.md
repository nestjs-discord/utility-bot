# Utility Bot

[![build-and-test](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml)
[![CodeFactor](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/badge/main)](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/overview/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/nestjs-discord/utility-bot)](https://goreportcard.com/report/github.com/nestjs-discord/utility-bot)

A Discord bot
designed to streamline the support process on [the official NestJS Discord server](https://discord.gg/nestjs).

We've developed this bot to improve your interaction with other users.
With pre-written, perfectly formatted Markdown content delivered through slick slash commands,
you can respond to everyday situations quickly.
No more tedious explanations, no more typing out the same thing over and over again.

## Commands

| Command                        | Description                                                                    | Who can execute it?    |
|--------------------------------|--------------------------------------------------------------------------------|------------------------|
| `/solved`                      | To close and mark the forum posts as solved with a tag                         | Mods, Forum post owner |
| `/archive`                     | To close and lock the forum posts                                              | Mods                   |
| `/dont-ping-mods`              | Tell someone to stop pinging mods for help                                     | Mods                   |
| `/mods cross-post`             | Warn users not to cross-post, ask in relevant channels, and be patient         | Mods                   |
| `/mods spam`                   | A friendly warning to prevent spamming in the chat                             | Mods                   |
| `/mods elaborate`              | Please provide more context so we can better assist you with your inquiry      | Everyone               |
| `/community awesome`           | A curated list of awesome things related to NestJS (`awesome-nestjs`)          | Everyone               |
| `/community testing`           | Show off to the community methods of testing NestJS (`jmcdo29/testing-nestjs`) | Everyone               |
| `/docs behind-proxy`           | How to enable the trust proxy option behind a reverse proxy                    | Everyone               |
| `/docs circular-dependency`    | Circular dependencies in NestJS and how to avoid them                          | Everyone               |
| `/docs file-change-detected`   | How to fix when TypeScript 4.9+ on infinite loop in watch mode                 | Everyone               |
| `/docs nest-debug`             | Debug NestJS dependencies with environment variable                            | Everyone               |
| `/docs providers-vs-imports`   | Ensure accurate service/module listings                                        | Everyone               |
| `/docs request-lifecycle`      | Get a grip on NestJS' request lifecycle for smooth coding                      | Everyone               |
| `/docs resolve-dependencies`   | Nest not being able to resolve dependencies of something                       | Everyone               |
| `/rules codeblocks`            | Letting people know how to share their code in a Markdown code block           | Everyone               |
| `/rules dm`                    | Use server to get help on NestJS; don't DM other members                       | Everyone               |
| `/rules dont-ask-to-ask`       | Telling people to just ask their question with a link explaining why           | Everyone               |
| `/rules no-hello`              | Telling people not to say just 'hello' in the chat                             | Everyone               |
| `/rules reproduction`          | Request for a minimum reproduction of the issue                                | Everyone               |
| `/rules screenshot`            | When someone posts a screenshot instead of sharing their code                  | Everyone               |
| `/javascript floating-promise` | Introduction to floating Promises in JavaScript                                | Everyone               |
| `/request-info`                | Ask people to run `@nestjs/cli info` and provide the output.                   | Everyone               |
| `/show-the-error`              | Telling people to be specific when seeking help.                               | Everyone               |

## Configuration

```shell
cp .env.example .env
```

The following environment variables are required, and the rest of the configuration is located in `config.yml`.

- `DISCORD_APP_ID`
- `DISCORD_BOT_TOKEN`
- `DISCORD_GUILD_ID`

## How we handle incoming Discord updates

```mermaid
flowchart TD
%% classes
    classDef orange stroke: #f96
    classDef green stroke: #0f0
%% start point
    discordUpdate["Discord Update"]
    discordUpdate -->|" Event "| utilityBotApp("Utility-Bot Application")
    utilityBotApp --> eventType{"What's the event type?"}
    eventType -->|" Ready "| setStatus["Set Bot status to:\n\nListening to slash-commands!"]
    setStatus --> finish["Finish"]
    eventType -->|" Interaction Create "| eventInteractionCreate{"What is the type of interaction?"}
    eventInteractionCreate -->|" Application Command Autocomplete "| applicationCommandAutoComplete["Slash command autocomplete flow"]:::orange
    eventInteractionCreate -->|" Application Command "| checkRateLimit("Check ratelimit policy")
    checkRateLimit --> isStaticCommand{"Is static command?"}
    isStaticCommand -->|" Yes "| handleStaticCommand{"Handle static command\n(hardcoded)"}
    handleStaticCommand -->|" /solved "| solvedCommandFlow["Solved cmd flow"]:::orange
    handleStaticCommand -->|" /archive "| archiveCommandFlow["Archive cmd flow"]:::orange
    handleStaticCommand -->|" /dont-ping-mods "| dontPingModsFlow["Don't ping mods cmd flow"]:::orange
    handleStaticCommand -->|" /google-it "| googleItFlow["Google it cmd flow"]:::orange
    handleStaticCommand -->|" /reference "| referenceFlow["Reference cmd flow"]:::orange
    isStaticCommand -->|" No "| handleDynamicCommand("Handle dynamic command\n(config.yml)"):::green
    eventType -->|Message Create| isBotMessage{"Is the message from a bot?"}
    isBotMessage -->|" No "| AutoModFlow["Auto-moderation feature flow"]:::orange
    isBotMessage -->|" Yes "| Finish["Finish"]
```

### Auto-moderation feature flow

We have noticed that a small number of members in our community have had
their Discord accounts taken over due to clicking on harmful links.
Attackers who gain access to these accounts often use them to send harmful links to multiple channels.
To prevent this, we have added an auto-moderation feature to our utility bot.
This feature will detect and instantly ban any accounts that are sending harmful links,
protecting our community members from clicking on them.
Even though the bot's primary function is to send pre-written responses through slash-commands,
we believe this additional feature will be invaluable in keeping our community safe.

```mermaid
flowchart TD
%% start point
    start["Start point"]
    start --> shouldTrackChannelID
    shouldTrackChannelID -->|" Yes "| isUserModerator
    isUserModerator{"Is message author\na moderator?"}
    isUserModerator -->|" Yes "| finish-1["Finish"]
    isUserModerator -->|" No "| isUserInDeniedList

    subgraph take-action-primary["Take action - primary"]
        deleteTheirMessages("Delete their message(s)")
        banThem("Ban their account")
        alertModerators("Alert moderators")
        
        deleteTheirMessages --> banThem
        banThem --> alertModerators
        alertModerators --> finish-5["Finish"]
    end

    subgraph take-action-secondary["Take action - secondary (in case the primary fails)"]
        deleteTheirMessage("Delete their current message")
        banThemAgain("Ban their account")
        deleteTheirMessage --> banThemAgain
        banThemAgain --> finish-3["Finish"]
    end

    subgraph In-memory temporary cache
%%        subgraph Databases
            channelsToTrack[("Channel IDs\nto moderate")]
            messages[("Discord messages\n\nmap[UserID]map[ChannelID]Message")]
            deniedList[("Denied users IDs")]
%%        end

    %% channels to track
        shouldTrackChannelID{"Should moderate the current\nchannel ID?\n(we only track text channels)"}
        shouldTrackChannelID -->|" No "| finish-2["Finish"]
        shouldTrackChannelID --- channelsToTrack
    %% denied list
        deniedList --- isUserInDeniedList
        isUserInDeniedList{"Is message author\nin the denied list?"}
        
        addUserToDeniedList("Add message author to the denied list")
        addUserToDeniedList --> deniedList

    %% Messages
        storeTheirMessage --> messages
        storeTheirMessage --> hasExceededMessageLimits
        hasExceededMessageLimits{"Has the message author\nexceeded the limits?"}
        hasExceededMessageLimits -->|" No "| finish-4["Finish"]
    end

    isUserInDeniedList -->|" Yes "| take-action-secondary
    isUserInDeniedList -->|" No "| storeTheirMessage("Store their message")
    hasExceededMessageLimits -->|" Yes "| take-action-primary-hub
    take-action-primary-hub(("."))
    take-action-primary-hub --> take-action-primary
    take-action-primary-hub --> addUserToDeniedList
```

### Solved command flow

Placeholder.

### Archive command flow

Placeholder.

### Don't ping moderators command flow

Placeholder.

### Google it command flow

Placeholder.

### Reference command flow

Placeholder.

### Slash command autocomplete flow

Placeholder.

## Docker usage

```shell
# brings down the previous container
# builds and starts a new container
# prints the container logs
make docker-redeploy
```
