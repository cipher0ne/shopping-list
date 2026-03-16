import { ShoppingList } from "./components/ShoppingList";
import styles from "./App.module.css";

function App() {
	return (
	<main>
		<h1 className={styles.title}>Shopping List</h1>
		<ShoppingList />
	</main>
	);
}

export default App;
