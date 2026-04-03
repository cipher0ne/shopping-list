import styles from "./styles/Menu.module.css";

type MenuProps = {
	selectedOption: string;
	setOption: (option: string) => void;
	setDisplayAuth: React.Dispatch<React.SetStateAction<boolean>>;
	userIsLoggedIn: boolean;
	setUserIsLoggedIn: React.Dispatch<React.SetStateAction<boolean>>;
}

export function Menu({
	selectedOption,
	setOption,
	setDisplayAuth,
	userIsLoggedIn,
	setUserIsLoggedIn
}: MenuProps) {
	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		setOption(e.target.value);
	}

	let userLabel = "";
	if (userIsLoggedIn) {
		const email = localStorage.getItem("email") || "";
		userLabel = email.split("@")[0];
	}

	const handleLoginLogout = () => {
		if (userIsLoggedIn) {
			setUserIsLoggedIn(false);
			localStorage.removeItem("email");
			localStorage.removeItem("token");
			localStorage.setItem("isRegistered", "true");
		} else {
			setDisplayAuth(true);
		}
	};

	return (
		<div className={styles.container}>
			<label>
				{userIsLoggedIn ? userLabel : null}
			</label>
			<button className={styles.loginButton} onClick={handleLoginLogout}>
				<span className={styles.userIcon}>{"\uf007"}</span>{" "}
				{userIsLoggedIn ? "Log out" : "Log in"}
			</button>
			<label className={styles.filterLabel}>{"\uf0b0"} Filter</label>
			<div className={styles.radioGroup}>
				<label>
					<input
						type="radio"
						name="filter"
						value="all"
						checked={selectedOption === "all"}
						onChange={handleChange}
					/>
					<span className={styles.radioDot}></span> Show all
				</label>
				<label>
					<input
						type="radio"
						name="filter"
						value={"not bought"}
						checked={selectedOption === "not bought"}
						onChange={handleChange}
					/>
					<span className={styles.radioDot}></span> Hide bought
				</label>
				<label>
					<input
						type="radio"
						name="filter"
						value="bought"
						checked={selectedOption === "bought"}
						onChange={handleChange}
					/>
					<span className={styles.radioDot}></span> Bought only
				</label>
			</div>
		</div>
	)
}
