import { createContext, useContext, useState} from "react";
import type { ReactNode } from "react";

interface AuthContextType {
    accessToken: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;
    login: (accessToken: string, refreshToken: string) => void;
}

const AuthContext = createContext<AuthContextType>({
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
    login: () => {},
});

function AuthProvider({children}: {children: ReactNode}) {
    
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const [refreshToken, setRefreshToken] = useState<string | null>(null);
    const isAuthenticated = accessToken !== null;
    const login = (accessToken: string, refreshToken: string) => {
        setAccessToken(accessToken);
        setRefreshToken(refreshToken);
    };

    return (
        <AuthContext.Provider 
            value={{
                accessToken,
                refreshToken,
                isAuthenticated,
                login,
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