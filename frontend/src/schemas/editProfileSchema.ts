import { z } from "zod";

export const editProfileSchema = z.object({
    username: z
        .string()
        .trim()
        .toLowerCase()
        .min(3, "Username must be at least 3 characters")
        .max(30, "Username must be at most 30 characters")
        .regex(/^[a-zA-Z0-9_-]+$/,"Login ID can only contain letters, numbers, underscores, and hyphens"),
});

export type EditProfileFormData = z.infer<typeof editProfileSchema>;