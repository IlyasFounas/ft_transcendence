import { z } from "zod";

export const signupSchema = z.object({
    email: z
        .string()
        .email("Please enter a valid email")
        .toLowerCase()
        .trim(),

    loginId: z
        .string()
        .trim()
        .toLowerCase()
        .min(3, "Login ID must be at least 3 characters")
        .max(20, "Login ID must be at most 30 characters")
        .regex(/^[a-zA-Z0-9_-]+$/,"Login ID can only contain letters, numbers, underscores, and hyphens"),

    username: z
        .string()
        .trim()
        .toLowerCase()
        .min(3, "Username must be at least 3 characters")
        .max(30, "Username must be at most 30 characters")
        .regex(/^[a-zA-Z0-9_-]+$/,"Username can only contain letters, numbers, underscores, and hyphens"),

    password: z
        .string()
        .min(8, "Password must be at least 8 characters")
        .max(50, "Password must be at most 50 characters")
        .refine(
            (value) => /[A-Z]/.test(value),
            "Password must contain at least one uppercase letter"
        )
        .refine(
            (value) => /[a-z]/.test(value),
            "Password must contain at least one lowercase letter"
        )
        .refine(
            (value) => /[0-9]/.test(value),
            "Password must contain at least one number"
        )
        .refine(
            (value) => /[^A-Za-z0-9]/.test(value),
            "Password must contain at least one special character"
        ),
});

export type SignupFormData = z.infer<typeof signupSchema>;