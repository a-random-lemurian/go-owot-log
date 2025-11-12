# Config
go-owot-log stores its config both in a configuration file and the environment
(with the prefix OWOTLOG). To set a configuration setting with environment
variables, uppercase all characters and replace dots with underscores. For
example, `chat.nickname: go-owot-log` would be `OWOTLOG_CHAT_NICKNAME=go-owot-log`.

## Root
`data_directory` is the directory in which go-owot-log stores its data. It should
not contain a trailing slash.

`worlds` is a list of worlds that the bot will log.

`global_receiver` is the world designated to process and receive global chat
messages. If this setting were not here, duplicates of global messages would
abound in the database.

`command_prefix` is the prefix that the bot uses for commands. This should not
contain any spaces.

## `ratelimits`
*For shorthand, assume that "ratelimits" is written before the first dot if a configuration
key begins with such.*
`ratelimits` contains settings pertaining to bot command ratelimits.

`.trusted_exempt` determines if trusted users will be exempt from ratelimits. To see
trusted users, refer to `users.trusted`.

`.commands` contains command-specific ratelimits. Each key under `.commands` should be the
name of a command, without the `ch` prefix. The values of these keys represent how often
the command may be run; for example, a value of `10s` means that the command may only be run
every ten seconds, or six times a minute.

## `chat`
`.nickname` is the nickname that the bot will use in chat.

`.color` is the color that the bot's nickname should be in. This will only work if the bot
is authenticated.

`.location` contains settings for the two chats of OWOT: global and page (world-specific).

### `chat.location.X`
`.enabled` is whether the bot should send messages in that chat. For example, it is
undesirable for both the production instance of the bot and staging instance to spam
global chat in response to user commands. `chat.location.global.enabled: false` fixes
this.

## `users`
`.trusted` is a list of trusted users. These users have administrator privileges. Entries
in `users` should use the Uvias or OWOT username of the trusted user.

## `lumberjack`
`lumberjack` contains settings for [lumberjack](https://github.com/natefinch/lumberjack),
a log rolling package that go-owot-log uses to write chat logs. The `lumberjack` object
contains the configuration of a lumberjack.Logger object. Refer to that repository for
up-to-date and more detailed information. For convenience, a copy is provided below.

`.filename` is the file to write logs to.

`.maxsize` is the maximum size of a log file before it is rotated.

`.maxage` is how long log files should be retained. Please remember that the bot has a
retention period of infinity.

`.maxbackups` is how many old log files should be retained.

`.localtime` is whether local time should be used instead of UTC.

`.compress` is whether old log files should be compressed. This often delivers storage
savings.
