/// <reference path="./global.d.ts" />
export declare class MicroLighter extends HTMLElement {
    #private;
    static observedAttributes: string[];
    constructor();
    connectedCallback(): void;
    disconnectedCallback(): void;
    attributeChangedCallback(): void;
}