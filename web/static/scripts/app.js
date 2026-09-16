import { CopyButton } from "./components/copyButton.js";
import { init as initCopyCodeBlock } from "./features/copyCodeBlock.js";
import { init as initFocusKey } from "./features/focusKey.js";
import { init as initGlobalCommands } from "./features/globalCommands.js";

CopyButton.define();

initCopyCodeBlock();
initFocusKey();
initGlobalCommands();
