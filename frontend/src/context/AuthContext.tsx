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
    logout: () => Promise<void>;
    refreshAccessToken: () => Promise<string | null>;
    authFetch: (url: string, options?: RequestInit) => Promise<{
        status: number,
        data: any,
    } | null >;
}

const AuthContext = createContext<AuthContextType>({
    accessToken: null,
    user: null,
    isAuthLoading: true,
    isAuthenticated: false,
    login: async () => {},
    logout: async () => {},
    refreshAccessToken: async () => null,
    authFetch: async () => null,
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

    const logout = async () => {
        
        const result = await authFetch("/api/auth/logout", {
            method: "POST",
        });

        if (result == null) {
            return;
        }
        if (result.status === 200) {
            setAccessToken(null);
            setUser(null);
        }
        return null;
    };

    const refreshAccessToken = async () => {

        const result = await apiFetch("/api/auth/refresh");

        if (result === null) {
            return null;
        }
        if (result.status === 200)
        {
            const newToken = result.data.Authorization;
            setAccessToken(newToken);
            
            await getUser(newToken);

            return newToken;
        }
        return null;
    }

    const authFetch = async (
        url: string,
        options: RequestInit = {}
    ) => {

        if (!accessToken) {
            return null;
        }

        const result = await apiFetch(url, options, accessToken);
        if (result === null) {
            return null;
        }
        if (result.status !== 401) {
            return result;
        }

        const newToken = await refreshAccessToken();
        if (newToken === null) {
            return null;
        }

        return await apiFetch(url, options, newToken);
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
                logout,
                refreshAccessToken,
                authFetch,
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