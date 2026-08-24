import { createContext, useContext, useState} from "react";
import type { ReactNode } from "react";

interface AuthContextType {
    authorizationToken: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;
    login: (authorizationToken: string, refreshToken: string) => void;
}

const AuthContext = createContext<AuthContextType>({
    authorizationToken: null,
    refreshToken: null,
    isAuthenticated: false,
    login: () => {},
});

function AuthProvider({children}: {children: ReactNode}) {
    
    const [authorizationToken, setAccessToken] = useState<string | null>(null);
    const [refreshToken, setRefreshToken] = useState<string | null>(null);
    const isAuthenticated = authorizationToken !== null;
    const login = (authorizationToken: string, refreshToken: string) => {
        setAccessToken(authorizationToken);
        setRefreshToken(refreshToken);
    };

    return (
        <AuthContext.Provider 
            value={{
                authorizationToken,
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