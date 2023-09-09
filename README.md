# ms2131x-status

This will get you the current status of an MS213X device in terms of input. It
will tell you whether it is currently receiving anything at all, the received
resolution and framerate.

Very much based on the work of [BertoldVdb](https://github.com/BertoldVdb/ms-tools)
whose libraries are doing the actual heavy lifting. The tool itself is mostly a
mangled version of his `cli` tool.

All thanks go to Bertold, and [markvdb](https://github.com/markvdb) for
[figuring out](https://github.com/BertoldVdb/ms-tools/issues/7#issuecomment-1706846783)
the memory locations! :)

Disclaimer: I don't really know what I'm doing, I'm just poking around and
trying to get the data I'm looking for.
