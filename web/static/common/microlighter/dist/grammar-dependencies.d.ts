/**
 * @import { Grammar, GrammarValue } from "./grammar.js"
 */
import type { Grammar, GrammarValue } from "./grammar.js";
export type LoadedGrammars = {
    /**
     * Grammars keyed by module name.
     */
    languages: Record<string, Grammar>;
    /**
     * Grammars keyed by TextMate scope name.
     */
    scopes: Map<string, Grammar>;
};
/**
 * Convert a supported language alias to its grammar module name.
 * @param {string} language
 * @param {Record<string, string>} [aliases]
 * @returns {string}
 */
export declare const normalizeLanguage: (language: string, aliases?: Record<string, string>) => string;
/**
 * Find grammar modules referenced by external TextMate scope includes.
 * Local repository references (`#...`, `$self`, and `$base`) are ignored.
 * @param {GrammarValue} value Grammar or nested grammar rule to inspect.
 * @returns {Set<string>} Normalized grammar module names.
 */
export declare const getExternalLanguages: (value: GrammarValue) => Set<string>;
/**
 * Create a cached loader that also resolves external grammar dependencies.
 * @param {(language: string) => Promise<Grammar | null>} [importLanguage]
 * @returns {(languages: string[]) => Promise<LoadedGrammars>}
 */
export declare const createGrammarLoader: (importLanguage?: (language: string) => Promise<Grammar | null>) => (languages: string[]) => Promise<LoadedGrammars>;
