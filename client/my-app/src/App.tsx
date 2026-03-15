import { useState } from "react"
import "./App.css"

type ProductT = { 
  id: number;
  name: string;
  quantity: number;
  bought: boolean;
};

type ProductHandlers = {
  increment: (id: number) => void;
  decrement: (id: number) => void;
}

function Product({p, handle}: {p: ProductT, handle: ProductHandlers}) {
  return (
    <section className="product">
      <div className="row">
        <input type="checkbox" name="myCheckBox" defaultChecked={p.bought} />
        <p>{p.name}</p>
      </div>

      <div className="row">
        <button className="button" onClick={() => handle.increment(p.id)}>
          +
        </button>
        
        <p>{p.quantity}</p>
        
        <button className="button" onClick={() => handle.decrement(p.id)}>
          -
        </button>
      </div>
    </section>
  )
}

function App() {
  const [products, setItems] = useState<ProductT[]>([]);

  const addItem = () => {
    setItems(prev => [
      { 
        id: prev.length + 1,
        name: "Name",
        quantity: 1,
        bought: false
      },
      ...prev
    ]);
  };

  const incrementQuantity = (id: number) => {
    setItems(prev => 
      prev.map(item =>
        item.id === id ? { ...item, quantity: item.quantity + 1} : item
      )
    );
  };

  const decrementQuantity = (id: number) => {
    setItems(prev => 
      prev.map(item => 
        item.id === id && item.quantity > 0 ? { ...item, quantity: item.quantity - 1} : item
      )
    );
  };

  return (
    <>
      <section id="list-section">
        <div id="list">
          <button className="button" onClick={addItem}>
            <span className="accent">+</span> New item
          </button>
          {products.map((item) => (
            <Product
              key={item.id}
              p={item}
              handle={{
                increment: incrementQuantity,
                decrement: decrementQuantity
              }}
            />
          ))}
        </div>
      </section>
    </>
  )
}

export default App
