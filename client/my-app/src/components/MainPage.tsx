import { ShoppingList } from "./ShoppingList";
import { Menu } from "./Menu"
import styles from "./MainPage.module.css";
import React, { useState } from "react";

type MainPageProps = {
	setDisplayAuth: React.Dispatch<React.SetStateAction<boolean>>;
}

export function MainPage({ setDisplayAuth }: MainPageProps) {
	const [displayMenu, setDisplayMenu] = useState(false);

	const toggleMenu = () => {
		displayMenu ? setDisplayMenu(false) : setDisplayMenu(true);
	}

	const [selectedOption, setOption] = useState("all");

	let menuContent;
	if (displayMenu) {
		menuContent = (
			<Menu
				selectedOption={selectedOption}
				setOption={setOption}
				setDisplayAuth={setDisplayAuth}
			/>
		)
	} else {
		menuContent = null;
	}

	return (
		<main className={styles.container}>
			<div className={styles.menu}>
				<button className={styles.menuButton} onClick={() => toggleMenu()}>
					{"\udb80\udf5c"}
				</button>
				{menuContent}
			</div>
			<h1 className={styles.title}>Shopping List</h1>
			<ShoppingList filter={selectedOption} />
		</main>
	);
}
