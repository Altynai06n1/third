import React, { useState } from "react";
import API from "../api";
import { Link } from "react-router-dom";

export default function Register() {
    const [form, setForm] = useState({ username: "", password: "" });
    const [message, setMessage] = useState("");

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await API.post("/register", form);
            setMessage("User registered successfully!");
        } catch {
            setMessage("Registration failed");
        }
    };

    return (
        <div className="container mt-5" style={{ maxWidth: 400 }}>
            <h3 className="mb-4 text-center">Register</h3>
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
                <button type="submit" className="btn btn-primary w-100">
                    Register
                </button>
            </form>
            <p className="text-center mt-3">
                Already have an account? <Link to="/login">Login</Link>
            </p>
            <p className="text-center text-muted mt-2">{message}</p>
        </div>
    );
}
