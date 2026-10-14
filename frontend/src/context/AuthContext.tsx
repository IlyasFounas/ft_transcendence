import { createContext, useContext, useState, useEffect, useRef } from "react";
import { apiFetch } from "../api/apiFetch";
import type { ReactNode } from "react";
import type { User } from "../types/user.types";

interface AuthContextType {
    user: User | null;
    isAuthLoading: boolean;
    isAuthenticated: boolean;
    login: () => Promise<void>;
    logout: () => Promise<void>;
    refreshAccessToken: () => Promise<boolean>;
    getUser: () => Promise<boolean>;
    authFetch: (url: string, options?: RequestInit) => Promise<{
        status: number,
        data: any,
    } | null >;
}

const AuthContext = createContext<AuthContextType>({
    user: null,
    isAuthLoading: true,
    isAuthenticated: false,
    login: async () => {},
    logout: async () => {},
    refreshAccessToken: async () => false,
    getUser: async () => false,
    authFetch: async () => null,
});

function AuthProvider({children}: {children: ReactNode}) {
    
    const [user, setUser] = useState<User | null>(null);
    const [isAuthLoading, setIsAuthLoading] = useState(true);
    const refreshPromiseRef = useRef<Promise<boolean> | null>(null);
    const isAuthenticated = user !== null;
    
    const refreshAccessToken = async () => {
        
        if (refreshPromiseRef.current) {
            return refreshPromiseRef.current;
        }

        const refreshPromise = (async () => {
            const result = await apiFetch("/api/auth/refresh");

            if (result === null) {
                return false;
            }

            return result.status === 200;
        })();

        refreshPromiseRef.current = refreshPromise;

        try {
            return await refreshPromise;
        } finally {
            refreshPromiseRef.current = null;
        }
    };
    
    const authFetch = async (
        url: string,
        options: RequestInit = {}
    ) => {

        const result = await apiFetch(url, options);
        if (result === null) {
            return null;
        }
        
        if (result.status !== 401) {
            return result;
        }

        const refreshed = await refreshAccessToken();
        if (!refreshed) {
            return null;
        }

        return await apiFetch(url, options);
    }

    const getUser = async (): Promise<boolean> => {
        const result = await apiFetch("/api/user");

        if (result === null) {
            setUser(null);
            return false;
        }

        if (result.status === 200) {
            setUser(result.data.user);
            return true;
        }

        if (result.status === 401) {
            setUser(null);
            return false;
        }

        setUser(null);
        return false;
    };

    
    const login = async () => {
        await getUser();
    };

    const logout = async () => {
        
        const result = await authFetch("/api/auth/logout", {
            method: "POST",
        });

        if (result == null) {
            return;
        }
        if (result.status === 200) {
            setUser(null);
        }
    };

    const authCheckStarted = useRef(false);

    useEffect(() => {
        if (authCheckStarted.current) return;
        authCheckStarted.current = true;

        const checkAuth = async () => {
            try {
                await getUser();
            } finally {
                setIsAuthLoading(false);
            }
        };

        void checkAuth();
    }, []);
    
    return (
        <AuthContext.Provider 
            value={{
                user,
                isAuthLoading,
                isAuthenticated,
                login,
                logout,
                refreshAccessToken,
                getUser,
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