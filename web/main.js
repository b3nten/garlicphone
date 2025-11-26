import { Renderable, html, css } from "./component.js";

class AppEntry extends Renderable {
	static styles = css`
		div {
			color: red;
		}
	`;
	onMounted() {
		console.log(this.shadowRoot.querySelector("div"))
	}
	render = () => html`<div>Hello, Garlic Phone!</div>`;
}

AppEntry.define()
