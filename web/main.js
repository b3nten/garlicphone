import { Renderable, html, css } from "./component.js";
import { Player } from "./messages.js";

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

fetch("/binary").then(x => x.arrayBuffer()).then(buffer => {
	console.log(new Uint8Array(buffer))
	console.log(new Player().fromBytes(new Uint8Array(buffer)))
});
