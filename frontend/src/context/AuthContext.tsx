import { createContext, useContext, useState} from "react";
import type { ReactNode } from "react";
import { apiFetch } from "../api/apiFetch";

interface AuthContextType {
    accessToken: string | null;
    isAuthenticated: boolean;
    login: (accessToken: string) => void;
    refreshAccessToken: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
    accessToken: null,
    isAuthenticated: false,
    login: () => {},
    refreshAccessToken: async () => {},
});

function AuthProvider({children}: {children: ReactNode}) {
    
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const isAuthenticated = accessToken !== null;
    const login = (accessToken: string) => {
        setAccessToken(accessToken);
    };

    const refreshAccessToken = async () => {
        const options: RequestInit = {
        method: "POST",
        }

        const result = await apiFetch("/api/auth/refresh", options);
        if (result === null) {
            return;
        }
        if (result.status === 200)
        {
            setAccessToken(result.data.Authorization);
        }
        return ;
    }
    return (
        <AuthContext.Provider 
            value={{
                accessToken,
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