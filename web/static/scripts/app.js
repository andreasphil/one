import { CopyButton } from "./components/copyButton.js";
import { StickyTitle } from "./components/stickyTitle.js";
import { init as initCopyCodeBlock } from "./features/copyCodeBlock.js";
import { init as initFocusKey } from "./features/focusKey.js";
import { init as initGlobalCommands } from "./features/globalCommands.js";
import { init as initHighlightCodeBlocks } from "./features/highlightCodeBlocks.js";
import { NoteView } from "./pages/note.js";
import { SearchView } from "./pages/search.js";

CopyButton.define();
StickyTitle.define();

NoteView.define();
SearchView.define();

initCopyCodeBlock();
initFocusKey();
initGlobalCommands();
initHighlightCodeBlocks();
