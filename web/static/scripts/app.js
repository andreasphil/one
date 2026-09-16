import { CopyButton } from "./components/copyButton.js";
import { init as initFocusKey } from "./features/focusKey.js";
import { init as initGlobalCommands } from "./features/globalCommands.js";

initFocusKey();
initGlobalCommands();

CopyButton.define();
