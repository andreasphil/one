//#region node_modules/.pnpm/lit-html@3.3.3/node_modules/lit-html/lit-html.js
/**
* @license
* Copyright 2017 Google LLC
* SPDX-License-Identifier: BSD-3-Clause
*/
var t$1 = globalThis;
var i$1 = (t) => t;
var s = t$1.trustedTypes;
var e$2 = s ? s.createPolicy("lit-html", { createHTML: (t) => t }) : void 0;
var h = "$lit$";
var o = `lit$${Math.random().toFixed(9).slice(2)}$`;
var n = "?" + o;
var r = `<${n}>`;
var l = document;
var c = () => l.createComment("");
var a = (t) => null === t || "object" != typeof t && "function" != typeof t;
var u = Array.isArray;
var d = (t) => u(t) || "function" == typeof t?.[Symbol.iterator];
var f = "[ 	\n\f\r]";
var v = /<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g;
var _ = /-->/g;
var m = />/g;
var p = RegExp(`>|${f}(?:([^\\s"'>=/]+)(${f}*=${f}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`, "g");
var g = /'/g;
var $ = /"/g;
var y = /^(?:script|style|textarea|title)$/i;
var x = (t) => (i, ...s) => ({
	_$litType$: t,
	strings: i,
	values: s
});
var b = x(1);
var E = Symbol.for("lit-noChange");
var A = Symbol.for("lit-nothing");
var C = /* @__PURE__ */ new WeakMap();
var P = l.createTreeWalker(l, 129);
function V(t, i) {
	if (!u(t) || !t.hasOwnProperty("raw")) throw Error("invalid template strings array");
	return void 0 !== e$2 ? e$2.createHTML(i) : i;
}
var N = (t, i) => {
	const s = t.length - 1, e = [];
	let n, l = 2 === i ? "<svg>" : 3 === i ? "<math>" : "", c = v;
	for (let i = 0; i < s; i++) {
		const s = t[i];
		let a, u, d = -1, f = 0;
		for (; f < s.length && (c.lastIndex = f, u = c.exec(s), null !== u);) f = c.lastIndex, c === v ? "!--" === u[1] ? c = _ : void 0 !== u[1] ? c = m : void 0 !== u[2] ? (y.test(u[2]) && (n = RegExp("</" + u[2], "g")), c = p) : void 0 !== u[3] && (c = p) : c === p ? ">" === u[0] ? (c = n ?? v, d = -1) : void 0 === u[1] ? d = -2 : (d = c.lastIndex - u[2].length, a = u[1], c = void 0 === u[3] ? p : "\"" === u[3] ? $ : g) : c === $ || c === g ? c = p : c === _ || c === m ? c = v : (c = p, n = void 0);
		const x = c === p && t[i + 1].startsWith("/>") ? " " : "";
		l += c === v ? s + r : d >= 0 ? (e.push(a), s.slice(0, d) + h + s.slice(d) + o + x) : s + o + (-2 === d ? i : x);
	}
	return [V(t, l + (t[s] || "<?>") + (2 === i ? "</svg>" : 3 === i ? "</math>" : "")), e];
};
var S = class S {
	constructor({ strings: t, _$litType$: i }, e) {
		let r;
		this.parts = [];
		let l = 0, a = 0;
		const u = t.length - 1, d = this.parts, [f, v] = N(t, i);
		if (this.el = S.createElement(f, e), P.currentNode = this.el.content, 2 === i || 3 === i) {
			const t = this.el.content.firstChild;
			t.replaceWith(...t.childNodes);
		}
		for (; null !== (r = P.nextNode()) && d.length < u;) {
			if (1 === r.nodeType) {
				if (r.hasAttributes()) for (const t of r.getAttributeNames()) if (t.endsWith(h)) {
					const i = v[a++], s = r.getAttribute(t).split(o), e = /([.?@])?(.*)/.exec(i);
					d.push({
						type: 1,
						index: l,
						name: e[2],
						strings: s,
						ctor: "." === e[1] ? I : "?" === e[1] ? L : "@" === e[1] ? z : H
					}), r.removeAttribute(t);
				} else t.startsWith(o) && (d.push({
					type: 6,
					index: l
				}), r.removeAttribute(t));
				if (y.test(r.tagName)) {
					const t = r.textContent.split(o), i = t.length - 1;
					if (i > 0) {
						r.textContent = s ? s.emptyScript : "";
						for (let s = 0; s < i; s++) r.append(t[s], c()), P.nextNode(), d.push({
							type: 2,
							index: ++l
						});
						r.append(t[i], c());
					}
				}
			} else if (8 === r.nodeType) if (r.data === n) d.push({
				type: 2,
				index: l
			});
			else {
				let t = -1;
				for (; -1 !== (t = r.data.indexOf(o, t + 1));) d.push({
					type: 7,
					index: l
				}), t += o.length - 1;
			}
			l++;
		}
	}
	static createElement(t, i) {
		const s = l.createElement("template");
		return s.innerHTML = t, s;
	}
};
function M(t, i, s = t, e) {
	if (i === E) return i;
	let h = void 0 !== e ? s._$Co?.[e] : s._$Cl;
	const o = a(i) ? void 0 : i._$litDirective$;
	return h?.constructor !== o && (h?._$AO?.(!1), void 0 === o ? h = void 0 : (h = new o(t), h._$AT(t, s, e)), void 0 !== e ? (s._$Co ??= [])[e] = h : s._$Cl = h), void 0 !== h && (i = M(t, h._$AS(t, i.values), h, e)), i;
}
var R = class {
	constructor(t, i) {
		this._$AV = [], this._$AN = void 0, this._$AD = t, this._$AM = i;
	}
	get parentNode() {
		return this._$AM.parentNode;
	}
	get _$AU() {
		return this._$AM._$AU;
	}
	u(t) {
		const { el: { content: i }, parts: s } = this._$AD, e = (t?.creationScope ?? l).importNode(i, !0);
		P.currentNode = e;
		let h = P.nextNode(), o = 0, n = 0, r = s[0];
		for (; void 0 !== r;) {
			if (o === r.index) {
				let i;
				2 === r.type ? i = new k(h, h.nextSibling, this, t) : 1 === r.type ? i = new r.ctor(h, r.name, r.strings, this, t) : 6 === r.type && (i = new Z(h, this, t)), this._$AV.push(i), r = s[++n];
			}
			o !== r?.index && (h = P.nextNode(), o++);
		}
		return P.currentNode = l, e;
	}
	p(t) {
		let i = 0;
		for (const s of this._$AV) void 0 !== s && (void 0 !== s.strings ? (s._$AI(t, s, i), i += s.strings.length - 2) : s._$AI(t[i])), i++;
	}
};
var k = class k {
	get _$AU() {
		return this._$AM?._$AU ?? this._$Cv;
	}
	constructor(t, i, s, e) {
		this.type = 2, this._$AH = A, this._$AN = void 0, this._$AA = t, this._$AB = i, this._$AM = s, this.options = e, this._$Cv = e?.isConnected ?? !0;
	}
	get parentNode() {
		let t = this._$AA.parentNode;
		const i = this._$AM;
		return void 0 !== i && 11 === t?.nodeType && (t = i.parentNode), t;
	}
	get startNode() {
		return this._$AA;
	}
	get endNode() {
		return this._$AB;
	}
	_$AI(t, i = this) {
		t = M(this, t, i), a(t) ? t === A || null == t || "" === t ? (this._$AH !== A && this._$AR(), this._$AH = A) : t !== this._$AH && t !== E && this._(t) : void 0 !== t._$litType$ ? this.$(t) : void 0 !== t.nodeType ? this.T(t) : d(t) ? this.k(t) : this._(t);
	}
	O(t) {
		return this._$AA.parentNode.insertBefore(t, this._$AB);
	}
	T(t) {
		this._$AH !== t && (this._$AR(), this._$AH = this.O(t));
	}
	_(t) {
		this._$AH !== A && a(this._$AH) ? this._$AA.nextSibling.data = t : this.T(l.createTextNode(t)), this._$AH = t;
	}
	$(t) {
		const { values: i, _$litType$: s } = t, e = "number" == typeof s ? this._$AC(t) : (void 0 === s.el && (s.el = S.createElement(V(s.h, s.h[0]), this.options)), s);
		if (this._$AH?._$AD === e) this._$AH.p(i);
		else {
			const t = new R(e, this), s = t.u(this.options);
			t.p(i), this.T(s), this._$AH = t;
		}
	}
	_$AC(t) {
		let i = C.get(t.strings);
		return void 0 === i && C.set(t.strings, i = new S(t)), i;
	}
	k(t) {
		u(this._$AH) || (this._$AH = [], this._$AR());
		const i = this._$AH;
		let s, e = 0;
		for (const h of t) e === i.length ? i.push(s = new k(this.O(c()), this.O(c()), this, this.options)) : s = i[e], s._$AI(h), e++;
		e < i.length && (this._$AR(s && s._$AB.nextSibling, e), i.length = e);
	}
	_$AR(t = this._$AA.nextSibling, s) {
		for (this._$AP?.(!1, !0, s); t !== this._$AB;) {
			const s = i$1(t).nextSibling;
			i$1(t).remove(), t = s;
		}
	}
	setConnected(t) {
		void 0 === this._$AM && (this._$Cv = t, this._$AP?.(t));
	}
};
var H = class {
	get tagName() {
		return this.element.tagName;
	}
	get _$AU() {
		return this._$AM._$AU;
	}
	constructor(t, i, s, e, h) {
		this.type = 1, this._$AH = A, this._$AN = void 0, this.element = t, this.name = i, this._$AM = e, this.options = h, s.length > 2 || "" !== s[0] || "" !== s[1] ? (this._$AH = Array(s.length - 1).fill(/* @__PURE__ */ new String()), this.strings = s) : this._$AH = A;
	}
	_$AI(t, i = this, s, e) {
		const h = this.strings;
		let o = !1;
		if (void 0 === h) t = M(this, t, i, 0), o = !a(t) || t !== this._$AH && t !== E, o && (this._$AH = t);
		else {
			const e = t;
			let n, r;
			for (t = h[0], n = 0; n < h.length - 1; n++) r = M(this, e[s + n], i, n), r === E && (r = this._$AH[n]), o ||= !a(r) || r !== this._$AH[n], r === A ? t = A : t !== A && (t += (r ?? "") + h[n + 1]), this._$AH[n] = r;
		}
		o && !e && this.j(t);
	}
	j(t) {
		t === A ? this.element.removeAttribute(this.name) : this.element.setAttribute(this.name, t ?? "");
	}
};
var I = class extends H {
	constructor() {
		super(...arguments), this.type = 3;
	}
	j(t) {
		this.element[this.name] = t === A ? void 0 : t;
	}
};
var L = class extends H {
	constructor() {
		super(...arguments), this.type = 4;
	}
	j(t) {
		this.element.toggleAttribute(this.name, !!t && t !== A);
	}
};
var z = class extends H {
	constructor(t, i, s, e, h) {
		super(t, i, s, e, h), this.type = 5;
	}
	_$AI(t, i = this) {
		if ((t = M(this, t, i, 0) ?? A) === E) return;
		const s = this._$AH, e = t === A && s !== A || t.capture !== s.capture || t.once !== s.once || t.passive !== s.passive, h = t !== A && (s === A || e);
		e && this.element.removeEventListener(this.name, this, s), h && this.element.addEventListener(this.name, this, t), this._$AH = t;
	}
	handleEvent(t) {
		"function" == typeof this._$AH ? this._$AH.call(this.options?.host ?? this.element, t) : this._$AH.handleEvent(t);
	}
};
var Z = class {
	constructor(t, i, s) {
		this.element = t, this.type = 6, this._$AN = void 0, this._$AM = i, this.options = s;
	}
	get _$AU() {
		return this._$AM._$AU;
	}
	_$AI(t) {
		M(this, t);
	}
};
var B = t$1.litHtmlPolyfillSupport;
B?.(S, k), (t$1.litHtmlVersions ??= []).push("3.3.3");
var D = (t, i, s) => {
	const e = s?.renderBefore ?? i;
	let h = e._$litPart$;
	if (void 0 === h) {
		const t = s?.renderBefore ?? null;
		e._$litPart$ = h = new k(i.insertBefore(c(), t), t, void 0, s ?? {});
	}
	return h._$AI(t), h;
};
//#endregion
//#region node_modules/.pnpm/lit-html@3.3.3/node_modules/lit-html/directive.js
/**
* @license
* Copyright 2017 Google LLC
* SPDX-License-Identifier: BSD-3-Clause
*/
var t = {
	ATTRIBUTE: 1,
	CHILD: 2,
	PROPERTY: 3,
	BOOLEAN_ATTRIBUTE: 4,
	EVENT: 5,
	ELEMENT: 6
};
var e$1 = (t) => (...e) => ({
	_$litDirective$: t,
	values: e
});
var i = class {
	constructor(t) {}
	get _$AU() {
		return this._$AM._$AU;
	}
	_$AT(t, e, i) {
		this._$Ct = t, this._$AM = e, this._$Ci = i;
	}
	_$AS(t, e) {
		return this.update(t, e);
	}
	update(t, e) {
		return this.render(...e);
	}
};
//#endregion
//#region node_modules/.pnpm/lit-html@3.3.3/node_modules/lit-html/directives/class-map.js
/**
* @license
* Copyright 2018 Google LLC
* SPDX-License-Identifier: BSD-3-Clause
*/ var e = e$1(class extends i {
	constructor(t$2) {
		if (super(t$2), t$2.type !== t.ATTRIBUTE || "class" !== t$2.name || t$2.strings?.length > 2) throw Error("`classMap()` can only be used in the `class` attribute and must be the only part in the attribute.");
	}
	render(t) {
		return " " + Object.keys(t).filter((s) => t[s]).join(" ") + " ";
	}
	update(s, [i]) {
		if (void 0 === this.st) {
			this.st = /* @__PURE__ */ new Set(), void 0 !== s.strings && (this.nt = new Set(s.strings.join(" ").split(/\s/).filter((t) => "" !== t)));
			for (const t in i) i[t] && !this.nt?.has(t) && this.st.add(t);
			return this.render(i);
		}
		const r = s.element.classList;
		for (const t of this.st) t in i || (r.remove(t), this.st.delete(t));
		for (const t in i) {
			const s = !!i[t];
			s === this.st.has(t) || this.nt?.has(t) || (s ? (r.add(t), this.st.add(t)) : (r.remove(t), this.st.delete(t)));
		}
		return E;
	}
});
//#endregion
//#region node_modules/.pnpm/nanostores@1.5.3/node_modules/nanostores/atom/index.js
var listenerQueue = [];
var lqIndex = 0;
var batchSeen = null;
var QUEUE_ITEMS_PER_LISTENER = 4;
var nanostoresGlobal = globalThis.nanostoresGlobal ||= { epoch: 0 };
var drainQueue = () => {
	let thrown;
	let i;
	while (lqIndex < listenerQueue.length) {
		i = lqIndex;
		lqIndex += QUEUE_ITEMS_PER_LISTENER;
		try {
			listenerQueue[i](listenerQueue[i + 1].value, listenerQueue[i + 2], listenerQueue[i + 3]);
		} catch (e) {
			thrown = e;
		}
	}
	listenerQueue.length = lqIndex = 0;
	if (thrown) throw thrown;
};
var atom = /* @__NO_SIDE_EFFECTS__ */ (initialValue) => {
	let listeners = [];
	let $atom = {
		eq: Object.is,
		get() {
			if (!$atom.lc) $atom.listen(() => {})();
			return $atom.value;
		},
		init: initialValue,
		lc: 0,
		listen(listener) {
			$atom.lc = listeners.push(listener);
			return () => {
				for (let i = lqIndex; i < listenerQueue.length;) if (listenerQueue[i] === listener) listenerQueue.splice(i, QUEUE_ITEMS_PER_LISTENER);
				else i += QUEUE_ITEMS_PER_LISTENER;
				let index = listeners.indexOf(listener);
				if (~index) {
					listeners.splice(index, 1);
					if (!--$atom.lc) $atom.off();
				}
			};
		},
		notify(oldValue, changedKey) {
			nanostoresGlobal.epoch++;
			let runListenerQueue = !listenerQueue.length && !batchSeen;
			for (let listener of listeners) {
				if (batchSeen?.has(listener)) continue;
				batchSeen?.add(listener);
				listenerQueue.push(listener, $atom, oldValue, batchSeen ? void 0 : changedKey);
			}
			if (runListenerQueue) drainQueue();
		},
		off() {},
		set(newValue) {
			let oldValue = $atom.value;
			if (!$atom.eq(oldValue, newValue)) {
				$atom.value = newValue;
				$atom.notify(oldValue);
			}
		},
		subscribe(listener) {
			let unbind = $atom.listen(listener);
			listener($atom.value);
			return unbind;
		},
		value: initialValue
	};
	return $atom;
};
//#endregion
//#region node_modules/.pnpm/nanostores@1.5.3/node_modules/nanostores/lifecycle/index.js
var MOUNT = 5;
var UNMOUNT = 6;
var REVERT_MUTATION = 10;
var on = (object, listener, eventKey, mutateStore) => {
	object.events = object.events || {};
	if (!object.events[eventKey + REVERT_MUTATION]) object.events[eventKey + REVERT_MUTATION] = mutateStore((eventProps) => {
		object.events[eventKey].reduceRight((event, l) => (l(event), event), {
			shared: {},
			...eventProps
		});
	});
	object.events[eventKey] = object.events[eventKey] || [];
	object.events[eventKey].push(listener);
	return () => {
		let currentListeners = object.events[eventKey];
		let index = currentListeners.indexOf(listener);
		if (~index) {
			currentListeners.splice(index, 1);
			if (!currentListeners.length) {
				object.events[eventKey + REVERT_MUTATION]();
				delete object.events[eventKey + REVERT_MUTATION];
			}
		}
	};
};
var STORE_UNMOUNT_DELAY = 1e3;
var onMount = ($store, initialize) => {
	let listener = (payload) => {
		let destroy = initialize(payload);
		if (destroy) $store.events[UNMOUNT].push(destroy);
	};
	return on($store, listener, MOUNT, (runListeners) => {
		let originListen = $store.listen;
		$store.listen = (...args) => {
			if (!$store.lc && !$store.active) {
				$store.active = true;
				runListeners();
			}
			return originListen(...args);
		};
		let originOff = $store.off;
		$store.events[UNMOUNT] = [];
		$store.off = () => {
			originOff();
			setTimeout(() => {
				if ($store.active && !$store.lc) {
					$store.active = false;
					for (let destroy of $store.events[UNMOUNT]) destroy();
					$store.events[UNMOUNT] = [];
				}
			}, STORE_UNMOUNT_DELAY);
		};
		return () => {
			$store.listen = originListen;
			$store.off = originOff;
		};
	});
};
//#endregion
//#region node_modules/.pnpm/nanostores@1.5.3/node_modules/nanostores/computed/index.js
var computedStore = (stores, cb, batched) => {
	if (!Array.isArray(stores)) stores = [stores];
	let previousArgs;
	let currentEpoch;
	let set = () => {
		if (currentEpoch === nanostoresGlobal.epoch) return;
		currentEpoch = nanostoresGlobal.epoch;
		let args = stores.map(($store) => $store.get());
		if (!previousArgs?.every((arg, i) => stores[i].eq(arg, args[i]))) {
			previousArgs = args;
			let value = cb(...args);
			if (value && value.then && value.t) value.then((asyncValue) => {
				if (previousArgs === args) $computed.set(asyncValue);
			});
			else {
				$computed.set(value);
				currentEpoch = nanostoresGlobal.epoch;
			}
		}
	};
	let $computed = /* @__PURE__ */ atom();
	let get = $computed.get;
	$computed.get = () => {
		set();
		return get();
	};
	let timer;
	let run = batched ? () => {
		clearTimeout(timer);
		timer = setTimeout(set);
	} : set;
	onMount($computed, () => {
		let unbinds = stores.map(($store) => $store.listen(run));
		set();
		return () => {
			for (let unbind of unbinds) unbind();
		};
	});
	return $computed;
};
var computed = /* @__NO_SIDE_EFFECTS__ */ (stores, fn) => computedStore(stores, fn);
//#endregion
//#region node_modules/.pnpm/nanostores@1.5.3/node_modules/nanostores/effect/index.js
var effect = (stores, callback) => {
	if (!Array.isArray(stores)) stores = [stores];
	let unbinds = [];
	let lastRunUnbind;
	let run = () => {
		lastRunUnbind && lastRunUnbind();
		lastRunUnbind = callback(...stores.map((store) => store.get()));
	};
	unbinds = stores.map((store) => store.listen(run));
	run();
	return () => {
		unbinds.forEach((unbind) => unbind());
		lastRunUnbind && lastRunUnbind();
	};
};
//#endregion
//#region node_modules/.pnpm/nanostores@1.5.3/node_modules/nanostores/map/index.js
var map = /* @__NO_SIDE_EFFECTS__ */ (initial = {}) => {
	let $map = /* @__PURE__ */ atom(initial);
	$map.eqKey = Object.is;
	$map.setKey = function(key, value) {
		let oldMap = $map.value;
		if (typeof value === "undefined" && key in $map.value) {
			$map.value = { ...$map.value };
			delete $map.value[key];
			$map.notify(oldMap, key);
		} else if (!$map.eqKey($map.value[key], value, key)) {
			$map.value = {
				...$map.value,
				[key]: value
			};
			$map.notify(oldMap, key);
		}
	};
	return $map;
};
//#endregion
//#region src/commandBar.js
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
function renderSvgFromString(svg) {
	return new DOMParser().parseFromString(svg, "image/svg+xml").documentElement;
}
var CommandBar = class CommandBar extends HTMLElement {
	static tag = "command-bar";
	static define(tag = this.tag) {
		this.tag = tag;
		customElements.define(tag, this);
	}
	static get instance() {
		const instance = document.querySelectorAll(this.tag);
		if (instance[0] instanceof CommandBar) {
			if (instance.length > 1) console.warn("Found multiple CommandBars. Only the first is returned.");
			return instance[0];
		} else throw new Error("No CommandBar instance found.");
	}
	/**
	* @type {import("nanostores").MapStore<{
	*   commands: Command[];
	*   focusedResult: number;
	*   mostRecent: Command | null;
	*   open: boolean;
	*   query: string;
	* }>}
	*/
	#state = /* @__PURE__ */ map({
		commands: [],
		focusedResult: 0,
		mostRecent: null,
		open: false,
		query: ""
	});
	/** @type {import("nanostores").Atom<Record<string, Command>>} */
	#chords = /* @__PURE__ */ computed(this.#state, (state) => state.commands.reduce((all, current) => {
		if (current.chord) all[current.chord] = current;
		return all;
	}, {}));
	/** @type {import("nanostores").ReadableAtom<Command[]>} */
	#results = /* @__PURE__ */ computed([this.#state, this.#chords], (state, chords) => {
		if (!this.#state.get().query) return [];
		const matchingChord = chords[state.query];
		const queryTokens = state.query.toLowerCase().split(" ");
		const results = state.commands.filter((i) => {
			if (matchingChord && matchingChord.id === i.id) return false;
			const commandStr = [
				i.name,
				...i.alias ?? [],
				i.groupName ?? ""
			].join(" ").toLowerCase();
			return queryTokens.every((token) => commandStr.includes(token));
		}).slice(0, 10).sort((a, b) => (b.weight ?? 0) - (a.weight ?? 0));
		if (matchingChord) results.unshift(matchingChord);
		return results;
	});
	/** @type {KeyboardShortcut} */
	shortcut = {
		key: "k",
		metaKey: true
	};
	allowRepeat = true;
	get placeholder() {
		return this.getAttribute("placeholder") ?? "Search...";
	}
	get emptyMessage() {
		return this.getAttribute("emptymessage") ?? "Sorry, couldnʼt find anything.";
	}
	/** @param {Command[]} toRegister */
	registerCommand(...toRegister) {
		const ids = toRegister.map((c) => c.id);
		this.removeCommand(...ids);
		const next = this.#state.get().commands.concat(...toRegister);
		this.#state.setKey("commands", next);
		return () => {
			this.removeCommand(...ids);
		};
	}
	/** @param {string[]} toRemove */
	removeCommand(...toRemove) {
		const next = this.#state.get().commands.filter((c) => !toRemove.includes(c.id));
		this.#state.setKey("commands", next);
		const mostRecent = this.#state.get().mostRecent;
		if (mostRecent && toRemove.includes(mostRecent.id)) this.#state.setKey("mostRecent", null);
	}
	open(initialQuery = "") {
		this.#state.setKey("query", initialQuery);
		this.#toggle(true);
	}
	#disconnectedController = new AbortController();
	constructor() {
		super();
	}
	connectedCallback() {
		this.#disconnectedController = new AbortController();
		addEventListener("keydown", (e) => this.#onToggleShortcut(e), { signal: this.#disconnectedController.signal });
		addEventListener("keydown", (e) => this.#onGlobalKeydown(e), { signal: this.#disconnectedController.signal });
		this.addEventListener("command", (e) => this.#onCommand(e), { signal: this.#disconnectedController.signal });
		effect([this.#state, this.#results], () => {
			this.#render();
		});
		this.#render();
	}
	disconnectedCallback() {
		this.#disconnectedController.abort();
	}
	get #dialog() {
		return this.querySelector("dialog");
	}
	#toggle(open = !this.#state.get().open) {
		if (open === this.#state.get().open) return;
		else if (open) {
			this.#dialog?.showModal();
			this.#state.setKey("open", true);
		} else {
			this.#dialog?.close();
			this.#onDialogClose();
		}
	}
	#onDialogClose() {
		if (!this.#state.get().open) return;
		this.#state.setKey("focusedResult", 0);
		this.#state.setKey("open", false);
		this.#state.setKey("query", "");
	}
	/** @param {KeyboardEvent} event */
	#onToggleShortcut(event) {
		const strictShortcut = Object.assign({
			altKey: false,
			ctrlKey: false,
			metaKey: false,
			shiftKey: false,
			key: ""
		}, this.shortcut);
		if (!Object.entries(strictShortcut).reduce((match, [key, value]) => {
			return match && event[key] === value;
		}, true)) return;
		this.#toggle();
		event.preventDefault();
		event.stopPropagation();
	}
	/** @param {CommandEvent} event */
	#onCommand(event) {
		if (event.command === "--open") this.open();
	}
	/** @param {KeyboardEvent} event */
	#onGlobalKeydown(event) {
		if (event.key === "." && event.metaKey) {
			this.#runMostRecent();
			event.preventDefault();
		}
		if (!this.#state.get().open) return;
		let cancel = true;
		if (event.key === "Escape") this.#onEsc();
		else if (event.key === "ArrowUp") this.#moveFocusUp();
		else if (event.key === "ArrowDown") this.#moveFocusDown();
		else if (event.key === "Enter") this.#runFocusedCommand();
		else cancel = false;
		if (cancel) event.preventDefault();
	}
	#onEsc() {
		if (this.#state.get().query) this.#state.setKey("query", "");
		else this.#toggle(false);
	}
	/** @param {KeyboardEvent} event */
	#onQueryChange(event) {
		if (!(event.target instanceof HTMLInputElement)) return;
		if (!this.#state.get().open) {
			event.target.value = this.#state.get().query;
			return;
		}
		this.#state.setKey("query", event.target.value);
	}
	#moveFocusDown() {
		const commandCount = this.#results.get().length;
		if (commandCount === 0) return;
		const next = this.#state.get().focusedResult + 1;
		this.#state.setKey("focusedResult", Math.min(commandCount - 1, next));
	}
	#moveFocusUp() {
		const commandCount = this.#results.get().length;
		if (commandCount === 0) return;
		const next = this.#state.get().focusedResult - 1;
		this.#state.setKey("focusedResult", next < 0 ? commandCount - 1 : next);
	}
	#runFocusedCommand() {
		const focused = this.#results.get().at(this.#state.get().focusedResult);
		if (focused) this.#runCommand(focused);
	}
	#runMostRecent() {
		const command = this.#state.get().mostRecent;
		if (command && this.allowRepeat) this.#runCommand(command);
	}
	/** @param {Command} command */
	#runCommand(command) {
		command.action();
		this.#state.setKey("mostRecent", command);
		this.#toggle(false);
	}
	/**
	* @param {object} state
	* @param {Command[]} state.results
	* @param {number} state.focusedResult
	* @param {string} state.emptyMessage
	* @param {string} state.query
	* @param {string} state.searchLabel
	*/
	#template = (state) => b`<dialog @close="${() => this.#onDialogClose()}">
      <style>
        @scope {
          :scope {
            --internal-padding: 1rem;

            &:has(input:invalid) output {
              display: none;
            }
          }

          header {
            font-weight: normal;
            margin: 0;
          }

          output {
            all: unset;
            max-height: calc(80dvh - 12rem);
            overflow: auto;

            > * {
              margin-top: 0.75rem;
            }
          }

          ul {
            display: flex;
            flex-direction: column;
            gap: 0.125rem;
            list-style-type: none;
            margin: 0.75rem 0 0 0;
            padding: 0;
          }

          li button {
            background: transparent;
            color: var(--c-fg);
            justify-content: start;
            outline-offset: var(--outline-inset);
            text-align: left;
            width: 100%;

            &:hover {
              background: var(--c-surface-variant-bg);
            }

            &.focused {
              background: var(--c-surface-variant-bg);
            }

            &.chord-match {
              background: light-dark(var(--primary-50), var(--primary-100));
              color: var(--primary-500);

              &:hover {
                background: light-dark(var(--primary-100), var(--primary-200));
              }

              .group-name {
                color: var(--primary-400);
              }

              .chord {
                border-color: var(--primary-200);
                color: var(--primary-400);
              }
            }

            > :empty {
              display: none;
            }
          }

          .group-name {
            color: var(--c-fg-variant);
            display: inline-block;
            flex: none;
            font-weight: var(--font-weight-normal);

            &:after {
              content: "›";
              margin-left: 0.75ch;
            }
          }

          .chord {
            background: var(--c-surface-bg);
            border-radius: var(--border-radius-small);
            border: var(--border-width) solid var(--c-border);
            color: var(--c-fg-variant);
            font-family: var(--font-family-mono);
            font-size: var(--font-size-mono);
            margin-left: auto;
            padding: 0 0.25rem;
          }
        }
      </style>

      <header>
        <label>
          <span class="visually-hidden">${state.searchLabel}</span>
          <input
            type="search"
            .value="${state.query}"
            placeholder="${state.searchLabel}"
            required
            @input="${(event) => this.#onQueryChange(event)}"
          />
        </label>
      </header>

      <output has-fallback="${state.results.length ? "" : "empty"}">
        <ul>
          ${state.results.map((r, i) => b`<li>
                <button
                  class="${e({
		focused: i === state.focusedResult,
		"chord-match": r.chord === state.query
	})}""
                  @click="${() => this.#runCommand(r)}"
                >
                  <span class="icon">${r.icon}</span>
                  <span class="group-name">${r.groupName}</span>
                  <span class="clamp">${r.name}</span>
                  <span class="chord">${r.chord}</span>
                </button>
              </li>`)}
        </ul>

        <div fallback-for="empty">
          <p>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="lucide lucide-frown"
            >
              <circle cx="12" cy="12" r="10" />
              <path d="M16 16s-1.5-2-4-2-4 2-4 2" />
              <line x1="9" x2="9.01" y1="9" y2="9" />
              <line x1="15" x2="15.01" y1="9" y2="9" />
            </svg>
          </p>
          <p data-test-id="empty-message">${state.emptyMessage}</p>
        </div>
      </output>
    </dialog>`;
	#render() {
		D(this.#template({
			emptyMessage: this.emptyMessage,
			focusedResult: this.#state.get().focusedResult,
			query: this.#state.get().query,
			results: this.#results.get(),
			searchLabel: this.placeholder
		}), this);
	}
};
//#endregion
export { CommandBar, renderSvgFromString };
