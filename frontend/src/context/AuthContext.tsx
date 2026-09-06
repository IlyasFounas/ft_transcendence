import { createContext, useContext, useState, useEffect } from "react";
import type { ReactNode } from "react";
import type { User } from "../types/user.types";
import { apiFetch } from "../api/apiFetch";

interface AuthContextType {
    accessToken: string | null;
    user: User | null;
    isAuthLoading: boolean;
    isAuthenticated: boolean;
    login: (accessToken: string) => Promise<void>;
    refreshAccessToken: () => Promise<boolean>;
}

const AuthContext = createContext<AuthContextType>({
    accessToken: null,
    user: null,
    isAuthLoading: true,
    isAuthenticated: false,
    login: async () => {},
    refreshAccessToken: async () => false,
});

function AuthProvider({children}: {children: ReactNode}) {
    
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const [user, setUser] = useState<User | null>(null);
    const [isAuthLoading, setIsAuthLoading] = useState(true);
    const isAuthenticated = accessToken !== null;
    
    const getUser = async (token: string) => {
        const result = await apiFetch("/api/user", {}, token);
        
        if (result === null) {
            return false;
        }
        if (result.status === 200)
        {
            setUser(result.data.user);
            return true;
        }
        return false;
        
    }
    
    const login = async (accessToken: string) => {
        setAccessToken(accessToken);
        await getUser(accessToken);
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
            const newToken = result.data.Authorization;
            setAccessToken(newToken);
            
            await getUser(newToken);

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
                user,
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