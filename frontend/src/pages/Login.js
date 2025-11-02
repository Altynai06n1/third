import React, { useState, useContext } from "react";
import { AuthContext } from "../AuthContext";
import API from "../api";
import { useNavigate, Link } from "react-router-dom";

export default function Login() {
    const [form, setForm] = useState({ username: "", password: "" });
    const [message, setMessage] = useState("");
    const { login } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const res = await API.post("/login", form);
            login(res.data.token);
            navigate("/items");
        } catch {
            setMessage("Invalid credentials");
        }
    };

    return (
        <div className="container mt-5" style={{ maxWidth: 400 }}>
            <h3 className="mb-4 text-center">Login</h3>
            <form onSubmit={handleSubmit}>
                <div className="mb-3">
                    <input
                        className="form-control"
                        placeholder="Username"
                        value={form.username}
                        onChange={(e) =>
                            setForm({ ...form, username: e.target.value })
                        }
                    />
                </div>
                <div className="mb-3">
                    <input
                        className="form-control"
                        placeholder="Password"
                        type="password"
                        value={form.password}
                        onChange={(e) =>
                            setForm({ ...form, password: e.target.value })
                        }
                    />
                </div>
                <button type="submit" className="btn btn-success w-100">
                    Login
                </button>
            </form>
            <p className="text-center mt-3">
                Don't have an account? <Link to="/register">Register</Link>
            </p>
            <p className="text-center text-muted mt-2">{message}</p>
        </div>
    );
}
