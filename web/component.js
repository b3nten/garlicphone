import { LitElement } from "lit";
export { html, css } from "lit";

const toKebab = (str) =>
	str
		.replace(/([a-z])([A-Z])/g, "$1-$2") // get all lowercase letters that are near to uppercase ones
		.replace(/[\s_]+/g, "-") // replace all spaces and low dash
		.toLowerCase(); // convert to lower case

export class Renderable extends LitElement {
	static define(maybeTag) {
		const tag = maybeTag || toKebab(this.name);
		queueMicrotask(() => {
			if (!customElements.get(tag)) {
				customElements.define(tag, this);
			}
		});
	}

	onMount() { }
	onMounted() { }
	onUpdate() { }
	onUpdated() { }
	onUnmount() { }

	connectedCallback() {
		super.connectedCallback();
		this.onMount();
	}

	firstUpdated() {
		this.onMounted()
	}

	willUpdate() {
		this.onUpdate()
	}

	updated() {
		this.onUpdated()
	}

	disconnectedCallback() {
		super.disconnectedCallback();
		this.onUnmount();
	}
};
