import React, { useState, useEffect, useContext } from "react";
import API from "../api";
import { AuthContext } from "../AuthContext";
import { useNavigate } from "react-router-dom";

export default function Items() {
    const [items, setItems] = useState([]);
    const [form, setForm] = useState({ name: "", description: "" });
    const [loading, setLoading] = useState(true); // ← ДОБАВИЛ
    const { logout, role } = useContext(AuthContext);
    const navigate = useNavigate();

    useEffect(() => {
        loadItems();
    }, []);

    const loadItems = async () => {
        setLoading(true);
        try {
            const res = await API.get("/api/items");
            setItems(res.data);
        } catch (err) {
            console.error("Load items error:", err);
            if (err.response?.status === 401) {
                logout();
                navigate("/login");
            }
        } finally {
            setLoading(false);
        }
    };

    const addItem = async (e) => {
        e.preventDefault();
        if (!form.name.trim() || !form.description.trim()) {
            alert("Заполните все поля");
            return;
        }
        try {
            await API.post("/api/items", {
                name: form.name.trim(),
                description: form.description.trim(),
            });
            setForm({ name: "", description: "" });
            loadItems();
        } catch (err) {
            console.error("Add item error:", err);
            if (err.response?.status === 403) {
                alert("Только админ может добавлять элементы");
            }
        }
    };

    const deleteItem = async (id) => {
        if (!window.confirm("Удалить элемент?")) return;
        try {
            await API.delete(`/api/items/${id}`);
            loadItems();
        } catch (err) {
            console.error("Delete failed:", err);
        }
    };

    return (
        <div className="container mt-4">
            <div className="d-flex justify-content-between align-items-center mb-3">
                <h3>Items</h3>
                <button
                    className="btn btn-outline-danger btn-sm"
                    onClick={() => {
                        logout();
                        navigate("/login");
                    }}
                >
                    Logout
                </button>
            </div>

            {role === "admin" && (
                <form onSubmit={addItem} className="mb-4">
                    <div className="row g-2">
                        <div className="col-md-4">
                            <input
                                className="form-control"
                                placeholder="Name"
                                value={form.name}
                                onChange={(e) => setForm({ ...form, name: e.target.value })}
                            />
                        </div>
                        <div className="col-md-6">
                            <input
                                className="form-control"
                                placeholder="Description"
                                value={form.description}
                                onChange={(e) =>
                                    setForm({ ...form, description: e.target.value })
                                }
                            />
                        </div>
                        <div className="col-md-2">
                            <button type="submit" className="btn btn-primary w-100">
                                Add
                            </button>
                        </div>
                    </div>
                </form>
            )}

            {loading ? (
                <p className="text-center text-muted">Loading items...</p>
            ) : items.length === 0 ? (
                <p className="text-center">No items yet.</p>
            ) : (
                <div className="row">
                    {items.map((i) => (
                        <div className="col-md-4 mb-3" key={i.id}>
                            <div className="card shadow-sm">
                                <div className="card-body">
                                    <h5 className="card-title">{i.name}</h5>
                                    <p className="card-text text-muted">{i.description}</p>
                                    {role === "admin" && (
                                        <div className="btn-group" role="group">
                                            <button
                                                className="btn btn-sm btn-danger"
                                                onClick={() => deleteItem(i.id)}
                                            >
                                                Delete
                                            </button>
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}