import styles from "./styles/AuthPage.module.css"
import { useEffect, useState } from "react"

type AuthDisplayProps = {
	setDisplayAuth: React.Dispatch<React.SetStateAction<boolean>>;
	setUserIsLoggedIn: React.Dispatch<React.SetStateAction<boolean>>;
}

export function AuthPage({ setDisplayAuth, setUserIsLoggedIn }: AuthDisplayProps) {
	const [isRegistered, setRegistered] = useState(
		() => localStorage.getItem("isRegistered") === "true"
	);


	// API endpoint base URL
	const API_BASE = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

	async function registerUser(email: string, password: string): Promise<{ success: boolean; message?: string }> {
		try {
			const res = await fetch(`${API_BASE}/register`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ email, password })
			});
			if (res.ok) {
				return { success: true };
			} else {
				const data = await res.json().catch(() => ({}));
				return { success: false, message: data?.error || "Registration failed" };
			}
		} catch (err) {
			return { success: false, message: "Network error" };
		}
	}

	async function loginUser(email: string, password: string): Promise<{ success: boolean; token?: string; message?: string }> {
		try {
			const res = await fetch(`${API_BASE}/login`, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ email, password })
			});
			if (res.ok) {
				const data = await res.json();
				return { success: true, token: data.token };
			} else {
				const data = await res.json().catch(() => ({}));
				return { success: false, message: data?.error || "Login failed" };
			}
		} catch (err) {
			return { success: false, message: "Network error" };
		}
	}

	const handleRegister = async () => {
		const result = await registerUser(email, password);
		if (result.success) {
			localStorage.setItem("isRegistered", "true");
			setRegistered(true);
			alert("Registration successful! Please sign in.");
		} else {
			alert(result.message || "Registration failed");
		}
	}

	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [rePassword, setRePassword] = useState("");

	const [canConfirm, setCanConfirm] = useState(false);

	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;

	useEffect(() => {
		const validEmail 		= emailRegex.test(email);
		const passwordLength 	= password.length >= 12;
		const hasLowerCase		= /[a-z]/.test(password);
		const hasUpperCase		= /[A-Z]/.test(password);
		const hasNumber 		= /\d/.test(password);
		const hasSpecialChar 	= /[!@#$%^&*]/.test(password);
		const onlyAllowed 		= /^[A-Za-z0-9!@#$%^&*]+$/.test(password);

		if (
			validEmail &&
			passwordLength &&
			hasLowerCase &&
			hasUpperCase &&
			hasNumber &&
			hasSpecialChar &&
			onlyAllowed &&
			(isRegistered || password === rePassword)
		) {
			setCanConfirm(true);
		} else {
			setCanConfirm(false);
		}
	}, [email, password, rePassword, isRegistered]);

	const [showPassword, setShowPassword] = useState(false);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	let title;
	let repeatPassword;
	let confirmButton;

	if (isRegistered) {
		title = <h2 className={styles.title}>Sign in</h2>;
		repeatPassword = null;
		confirmButton = (
			<button
				className={styles.confirmButton}
				disabled={!canConfirm || loading}
				onClick={async () => {
					setLoading(true);
					setError(null);
					const result = await loginUser(email, password);
					setLoading(false);
					if (result.success && result.token) {
						localStorage.setItem("token", result.token);
						setUserIsLoggedIn(true);
						setDisplayAuth(false);
					} else {
						setError(result.message || "Login failed");
					}
				}}
			>
				{loading ? "Signing in..." : "Sign in"}
			</button>
		);
	} else {
		title = <h2 className={styles.title}>Sign up</h2>;
		repeatPassword = (
			<>
				<label className={styles.label}>repeat password</label>
				<input
					className={styles.inputField}
					type={showPassword ? "text" : "password"}
					onChange={(e) => setRePassword(e.target.value.trim())}
				/>
			</>
		);
		confirmButton = (
			<button
				className={styles.confirmButton}
				onClick={async () => {
					setLoading(true);
					setError(null);
					await handleRegister();
					setLoading(false);
				}}
				disabled={!canConfirm || loading}
			>
				{loading ? "Signing up..." : "Sign up"}
			</button>
		);
	}

	return (
		<div className={styles.container}>
			<section className={styles.authForm}>
				{title}

				<label className={styles.label}>e-mail</label>
				<input
					className={styles.inputField}
					type="text"
					onChange={(e) => setEmail(e.target.value)}
				/>

				<label className={styles.label}>password</label>
				<div className={styles.passwordRow}>
					<input
						className={styles.inputField}
						type={showPassword ? "text" : "password"}
						onChange={(e) => setPassword(e.target.value.trim())}
					/>
					<button
						className={styles.eyeButton}
						onClick={() => {
							setShowPassword((v) => !v);
						}}
					>
						{showPassword ? "\udb80\ude09" : "\udb80\ude08"}
					</button>
				</div>
				{repeatPassword}

				{error && <div className={styles.error}>{error}</div>}

				<div className={styles.buttonRow}>
					<button className={styles.cancelButton} onClick={() => setDisplayAuth(false)}>
						Cancel
					</button>
					{confirmButton}
				</div>
			</section>
		</div>
	);
}
