import PrimaryButton from "../buttons/PrimaryButton";

export default function AuthenticationWelcome () {
    return (
        <div className="
            flex
            flex-col
            gap-[50px]
            py-[50px]
            "
        >
            <div className="flex flex-col gap-[30px]">
                <button className="flex h-[62px] bg-white rounded-2xl text-black text-heading-medium items-center justify-center">Continue with Google</button>
                <button className="flex h-[62px] bg-primary rounded-2xl text-black text-heading-medium items-center justify-center">Continue with Email</button>
            </div>
            <div className="flex gap-[16px] items-center">
                <hr className="w-full border-t border-text-secondary"></hr>
                <span className="text-body-primary text-text-secondary uppercase">or</span>
                <hr className="w-full border-t border-text-secondary"></hr>
            </div>
            <div className="flex items-center justify-between">
                <div className="text-heading-small text-text-primary">Already have an account ?</div>
                <PrimaryButton size="small">Login</PrimaryButton>
            </div>
        </div>
    );
}