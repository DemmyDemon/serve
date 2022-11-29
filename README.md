# serve
Quickly make a directory available over HTTP

## Why?
I needed a way to just open up a directory for file download temporarily. With this, I can just `serve -allow $MY_HOME_IP` to fetch whatever file from the current directory. Done.

## How?

There are not a lot of options, so let's just do this.

**-allow** lets you allow access from an IP. You need at least one of these, and wildcards are not supported.  
*Example:* `serve -allow 192.168.1.7 -allow 127.0.0.1`

Incidentally, yes, you even need to `-allow` the loopback address. There is no default value.

**-port** lets you specify what port you want to open. **The default port is 8181**.  
*Example:* `serve -allow 192.168.1.7 -port 1024`

Anything left over after *-allow* and *-port* is considered a filename to be served.  
The whole point of `serve` is to serve files in the current working directory, so paths make little sense here, and will be stripped. If you need to serve a file in a different directory, make that directory your current working directory.  
**The default is to serve all files in the current directory.**

## How do I run it as a daemon?

You don't. Seriously. Don't.

## Where is the settings file?

There isn't one.

## HTTPS?

No. You don't want that. If you need encryption, use SSH/SFTP/SCP.  
Baking in the complexity of certificates your browser will accept? Not worth it.

# Future?
I don't know. It might need some bugfixing, but I can't think of any more useful features.

Maaaaaybe look into detecting if the current terminal is in an SSH session, and if it is, suss out where it's coming from, and automagically doing an `-allow` for it?  
Maybe.
