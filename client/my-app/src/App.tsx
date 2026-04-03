import { MainPage } from "./components/MainPage";
import { AuthPage } from "./components/AuthPage";
import { useState } from "react";

function App() {
	const [displayAuth, setDisplayAuth] = useState(false);
	const [userIsLoggedIn, setUserIsLoggedIn] = useState(
		() => !!localStorage.getItem("token")
	);

	var content;
	if (displayAuth) {
		content = 	<AuthPage
						setDisplayAuth={setDisplayAuth}
						setUserIsLoggedIn={setUserIsLoggedIn}
					/>
	} else {
		content = 	<MainPage
						setDisplayAuth={setDisplayAuth}
						userIsLoggedIn={userIsLoggedIn}
						setUserIsLoggedIn={setUserIsLoggedIn}
					/>
	}

	return (<>{content}</>);
}

export default App;
