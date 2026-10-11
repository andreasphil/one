import type { Grammar, GrammarRule } from "./grammar.js";
export type HighlightAllOptions = {
    /**
     * Element or document to search.
     */
    root?: ParentNode;
    /**
     * Selector used to find code blocks.
     */
    selector?: string;
    /**
     * Extra aliases mapped to bundled grammar names.
     */
    languageAliases?: Record<string, string>;
};
export type MatchIndices = RegExpIndicesArray & {
    0: [number, number];
};
export type MatchOffsets = {
    indices: MatchIndices;
};
export type Match = RegExpExecArray & MatchOffsets;
export type RuleContext = {
    /**
     * Grammar the rule was resolved from.
     */
    grammar: Grammar;
    /**
     * Directly matchable rule.
     */
    rule: GrammarRule;
};
/**
 * Find code blocks, load their TextMate grammars, and register CSS highlights.
 * @param {HighlightAllOptions} [options]
 * @returns {Promise<HTMLElement[]>}
 */
export declare const highlightAll: ({ root, selector, languageAliases }?: HighlightAllOptions) => Promise<HTMLElement[]>;
