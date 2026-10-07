import { z } from "zod";

export const editProfileSchema = (currentUsername: string) => z.object({
    username: z
        .string()
        .trim()
        .toLowerCase()
        .min(3, "Username must be at least 3 characters")
        .max(30, "Username must be at most 30 characters")
        .refine(
            value =>
                value === currentUsername.trim().toLowerCase() ||
                /^[a-zA-Z0-9_-]+$/.test(value),
            {
                message:
                    "Login ID can only contain letters, numbers, underscores, and hyphens",
            }
        ),
    });

export type EditProfileFormData = {
    username: string;
};