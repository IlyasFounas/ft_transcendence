import { createContext, useContext, useState, useEffect } from "react";
import type { ReactNode } from "react";
import { apiFetch } from "../api/apiFetch";

interface AuthContextType {
    accessToken: string | null;
    isAuthLoading: boolean;
    isAuthenticated: boolean;
    login: (accessToken: string) => void;
    refreshAccessToken: () => Promise<boolean>;
}

const AuthContext = createContext<AuthContextType>({
    accessToken: null,
    isAuthLoading: true,
    isAuthenticated: false,
    login: () => {},
    refreshAccessToken: async () => false,
});

function AuthProvider({children}: {children: ReactNode}) {
    
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const [isAuthLoading, setIsAuthLoading] = useState(true);
    const isAuthenticated = accessToken !== null;
    const login = (accessToken: string) => {
        setAccessToken(accessToken);
    };

    const refreshAccessToken = async () => {
        const options: RequestInit = {
        method: "GET",
        }

        const result = await apiFetch("/api/auth/refresh", options);

        if (result === null) {
            return false;
        }
        if (result.status === 200)
        {
            setAccessToken(result.data.Authorization);
            return true;
        }
        return false;
    }
    
    useEffect(() => {
        const checkAuth = async () => {
            if (!accessToken) {
                await refreshAccessToken();
            }

            setIsAuthLoading(false);
        };

        checkAuth();
    }, []);
    
    return (
        <AuthContext.Provider 
            value={{
                accessToken,
                isAuthLoading,
                isAuthenticated,
                login,
                refreshAccessToken,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

function useAuth() {
    return useContext(AuthContext);
}

export { AuthProvider, useAuth };