import { base } from "./base.js";

export const nextjs = [
  ...base,
  {
    files: ["**/*.tsx"],
    rules: {
      "react/jsx-key": "off",
    },
  },
];

export default nextjs;
