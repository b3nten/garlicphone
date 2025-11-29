import { Renderable, html, css } from "./component.js";

class AppEntry extends Renderable {
	static styles = css`
		div {
			color: red;
		}
	`;
	render = () => html`<button @click=${() => console.log("LO")}>benchmark</div>`;
}

AppEntry.define()
