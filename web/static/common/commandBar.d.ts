export type Command = {
    /**
     * The unique identifier of the command. Can be any string.
     */
    id: string;
    /**
     * The visible name of the command.
     */
    name: string;
    /**
     * A list of aliases of the command. If the user searches for one of
     * them, the alias will be treated as if it was the name of the command.
     */
    alias?: string[];
    /**
     * A unique combination of characters. If the user types those exact
     * characters in the search field, the associated command will be shown prominently and
     * highlighted.
     */
    chord?: string;
    /**
     * An additional label displayed before the name.
     */
    groupName?: string;
    /**
     * Icon of the command. Should be a string (which will be
     * inserted as text content) or an HTML element (which will be inserted as-is).
     */
    icon?: string | HTMLElement;
    /**
     * Callback to run when the command is invoked.
     */
    action: () => void;
    /**
     * Used for sorting. Items with a higher weight will always appear
     * before items with a lower weight.
     */
    weight?: number;
};
export type KeyboardShortcut = Partial<Pick<KeyboardEvent, "key" | "metaKey" | "altKey" | "ctrlKey" | "shiftKey">>;
/**
 * @typedef Command
 * @property {string} id The unique identifier of the command. Can be any string.
 * @property {string} name The visible name of the command.
 * @property {string[]} [alias] A list of aliases of the command. If the user searches for one of
 *   them, the alias will be treated as if it was the name of the command.
 * @property {string} [chord] A unique combination of characters. If the user types those exact
 *   characters in the search field, the associated command will be shown prominently and
 *   highlighted.
 * @property {string} [groupName] An additional label displayed before the name.
 * @property {string | HTMLElement} [icon] Icon of the command. Should be a string (which will be
 *   inserted as text content) or an HTML element (which will be inserted as-is).
 * @property {() => void} action Callback to run when the command is invoked.
 * @property {number} [weight] Used for sorting. Items with a higher weight will always appear
 *   before items with a lower weight.
 */
/** @typedef {Partial<Pick<KeyboardEvent, "key" | "metaKey" | "altKey" | "ctrlKey" | "shiftKey">>} KeyboardShortcut */
/**
 * Takes an SVG string and converts it into an HTML element. Useful for displaying icons in the
 * command bar.
 *
 * @param {string} svg
 * @returns {HTMLElement}
 */
export declare function renderSvgFromString(svg: string): HTMLElement;
export declare class CommandBar extends HTMLElement {
    #private;
    static tag: string;
    static define(tag?: string): void;
    static get instance(): CommandBar;
    /** @type {KeyboardShortcut} */
    shortcut: KeyboardShortcut;
    allowRepeat: boolean;
    get placeholder(): string;
    get emptyMessage(): string;
    /** @param {Command[]} toRegister */
    registerCommand(...toRegister: Command[]): () => void;
    /** @param {string[]} toRemove */
    removeCommand(...toRemove: string[]): void;
    open(initialQuery?: string): void;
    constructor();
    connectedCallback(): void;
    disconnectedCallback(): void;
}
