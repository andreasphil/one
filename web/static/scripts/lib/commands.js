import { CommandBar } from "../../common/commandBar.js";

/** @typedef {import("../../common/commandBar.js").Command} Command */

/** @param {string} url */
export function newNavigationAction(url) {
  return () => {
    navigation.navigate(url);
  };
}

function createUseCommands() {
  CommandBar.define();

  return {
    register: CommandBar.instance.registerCommand.bind(CommandBar.instance),
  };
}

/** @type {ReturnType<typeof createUseCommands>} */
let _useCommands;

export function useCommands() {
  return (_useCommands ??= createUseCommands());
}
