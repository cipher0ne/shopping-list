import { useState } from "react";
import type { ProductT } from "../types";
import styles from "./ItemForm.module.css";

interface ItemFormProps {
	initialData?: ProductT;
	onConfirm: (name: string, quantity: number) => void;
	onCancel: () => void;
}

export function ItemForm({ initialData, onConfirm, onCancel }: ItemFormProps) {
	const [name, setName] = useState(initialData?.name || "");
	const [quantity, setQuantity] = useState(initialData?.quantity || 1);

	return (
		<div className={styles.formCard}>
			<input 
				className={styles.inputField}
				type="text" 
				value={name} 
				onChange={(e) => setName(e.target.value)} 
				placeholder="Product name"
				autoFocus
			/>
			<div className={styles.quantityRow}>
				<button className={styles.quantityButton} onClick={() => setQuantity(q => q + 1)}>+</button>
				<span>{quantity}</span>
				<button className={styles.quantityButton} onClick={() => setQuantity(q => Math.max(1, q - 1))}>-</button>
				<div className={styles.spacer} />
				<button className={styles.confirmButton} onClick={() => onConfirm(name, quantity)}>{"\udb80\udd2c"}</button>
				<button className={styles.cancelButton} onClick={onCancel}>{"\uf467"}</button>
				<div />
			</div>
		</div>
	);
}
