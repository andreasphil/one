import { CopyButton } from "./components/copyButton.js";
import { StickyTitle } from "./components/stickyTitle.js";
import { init as initCopyCodeBlock } from "./features/copyCodeBlock.js";
import { init as initFocusKey } from "./features/focusKey.js";
import { init as initGlobalCommands } from "./features/globalCommands.js";
import { init as initHighlightCodeBlocks } from "./features/highlightCodeBlocks.js";
import { init as initListNavigation } from "./features/listNavigation.js";
import { NoteView } from "./pages/note.js";

CopyButton.define();
StickyTitle.define();

NoteView.define();

initCopyCodeBlock();
initFocusKey();
initGlobalCommands();
initHighlightCodeBlocks();
initListNavigation("#note h2[id]");
initListNavigation(".search-results > li", { activate: ":scope > header a" });
