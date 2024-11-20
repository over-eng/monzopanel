import React from 'react';
import styles from "./button.module.css";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  asChild?: boolean;
  children: React.ReactNode;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ asChild = false, children, ...props }) => {
    const Comp = asChild ? 'span' : 'button';
    
    return (
      <Comp
        className={styles.button}
        {...props}
      >
        {children}
      </Comp>
    );
  }
);

Button.displayName = 'Button';

export default Button;