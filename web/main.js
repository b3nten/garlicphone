import { Renderable, html, css } from "./component.js";
import { deserialize, Player as Player2, Item as Item2 } from "./messages.js";
import { Item as Item1, Player as Player1 } from "./bebop.js"


class Instrumentor {
	static start(name) {
		performance.mark(`${name}::start`);
	}
	static end(name) {
		performance.mark(`${name}::end`);
		performance.measure(name, `${name}::start`, `${name}::end`);
	}
}

class AppEntry extends Renderable {
	static styles = css`
		div {
			color: red;
		}
	`;
	render = () => html`<button @click=${benchmark}>benchmark</div>`;
}

AppEntry.define()

function benchmark() {
	const player1 = Player1({
		id: 42,
		name: "Benton",
		// inventory: [],
		inventory: Array.from({length: 50}).map((_, i) => Item1({name: `Item ${i}`})),
		foo: "bar",
		dead: false
	})

	let bytes = Player1.encode(player1);

	let t = performance.now();
	for (let i = 0; i < 10000; i++) {
		let decoded = Player1.decode(bytes);
	}
	console.log("Bebop decode time:", performance.now() - t, Player1.decode(bytes));

	const player2 = new Player2();
	player2.id = 42;
	player2.name = "Benton";
	player2.inventory = Array.from({length: 50}).map((_, i) => {
		const item = new Item2();
		item.name = `Item ${i}`;
		return item;
	});
	player2.foo = "bar";
	player2.dead = false;

	bytes = player2.toBytes();

	t = performance.now();
	for (let i = 0; i < 10000; i++) {
		let decoded = deserialize(bytes)
	}
	console.log("GarlicPhone decode time:", performance.now() - t, deserialize(bytes));
}
