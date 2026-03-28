import { MainPage } from "./components/MainPage";
import { AuthPage } from "./components/AuthPage";
import { useState } from "react";

function App() {
	const [displayAuth, setDisplayAuth] = useState(false);

	return (
		<>
			{displayAuth ? <AuthPage setDisplayAuth={setDisplayAuth}/> : <MainPage setDisplayAuth={setDisplayAuth}/>}
		</>
	);
}

export default App;
