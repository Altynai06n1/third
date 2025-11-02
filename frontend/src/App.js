// App.jsx
import React, { useContext } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import { AuthContext } from "./AuthContext"; // ← импортируй
import Login from "./pages/Login";
import Register from "./pages/Register";
import Items from "./pages/Items";

// Встроенный PrivateRoute
function ProtectedRoute({ children }) {
    const { token } = useContext(AuthContext);
    return token ? children : <Navigate to="/login" replace />;
}

function App() {
    return (
        <Routes>
            <Route path="/" element={<Login />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route
                path="/items"
                element={
                    <ProtectedRoute>
                        <Items />
                    </ProtectedRoute>
                }
            />
        </Routes>
    );
}

export default App;