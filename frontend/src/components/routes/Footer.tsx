import { Link } from 'react-router-dom';

const Footer = () => {
  return (
	<footer className="flex justify-center gap-6 py-4 text-sm text-gray-400">
	  <Link to="/privacy" className="hover:text-white">Privacy Policy</Link>
	  <Link to="/terms" className="hover:text-white">Terms of Service</Link>
	</footer>
  );
};

export default Footer;