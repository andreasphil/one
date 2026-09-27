import { CopyButton } from "./components/copyButton.js";
import { StickyTitle } from "./components/stickyTitle.js";
import { init as initCopyCodeBlock } from "./features/copyCodeBlock.js";
import { init as initFocusKey } from "./features/focusKey.js";
import { init as initGlobalCommands } from "./features/globalCommands.js";
import { init as initSearchResultsNavigation } from "./features/searchResultsNavigation.js";

CopyButton.define();
StickyTitle.define();

initCopyCodeBlock();
initFocusKey();
initGlobalCommands();
initSearchResultsNavigation();
