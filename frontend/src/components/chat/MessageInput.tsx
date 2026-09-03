import { Send } from "lucide-react";
import { useRef } from "react";

interface MessageInputProps {
    value: string;
    onChange: (value: string) => void;
    onSend: () => void;
}

export default function MessageInput({
    value,
    onChange,
    onSend,
}: MessageInputProps) {

    const textareaRef = useRef<HTMLTextAreaElement>(null);

    const handleChange = (
        event: React.ChangeEvent<HTMLTextAreaElement>
    ) => {
        const textarea = event.target;
        onChange(textarea.value);

        textarea.style.height = "auto";
        textarea.style.height = `${Math.min(textarea.scrollHeight, 80)}px`;
    };

    const handleKeyDown = (
        event: React.KeyboardEvent<HTMLTextAreaElement>
    ) => {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();

            if (value.trim() !== "") {
                onSend();

                if (textareaRef.current) {
                    textareaRef.current.style.height = "auto";
                }
            }
        }
    };

    const handleSend = () => {
        if (value.trim() === "") {
            return;
        }

        onSend();

        if (textareaRef.current) {
            textareaRef.current.style.height = "auto";
        }
    };

    return (
        <div
            className="
                shrink-0
                border-t
                border-border-default
                p-3
            "
        >
            <div
                className="
                    flex
                    items-end
                    gap-2
                "
            >
                <textarea
                    ref={textareaRef}
                    value={value}
                    onChange={handleChange}
                    onKeyDown={handleKeyDown}
                    rows={1}
                    placeholder="Type a message..."
                    className="
                        min-w-0
                        max-h-20
                        flex-1
                        px-2.5
                        resize-none
                        overflow-y-auto
                        bg-surface-base
                        text-body-primary
                        text-text-primary
                        outline-none
                        placeholder:text-text-secondary
                        custom-scrollbar
                    "
                />

                <button
                    type="button"
                    onClick={handleSend}
                    className="
                        shrink-0
                        pb-1
                        text-text-secondary
                        transition
                        duration-150
                        hover:text-primary
                    "
                >
                    <Send size={20} />
                </button>
            </div>
        </div>
    );
}