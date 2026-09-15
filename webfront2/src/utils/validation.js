/**
 * Validation utilities for Office1789
 * Provides email and password strength validation
 */

/**
 * Validates email format
 * @param {string} email - Email to validate
 * @returns {boolean} True if valid
 */
export function isValidEmail(email) {
  if (!email) return false
  
  // RFC 5322 simplified regex for email validation
  const emailRegex = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/
  
  return emailRegex.test(email)
}

/**
 * Validates password strength
 * Requirements:
 * - At least 8 characters
 * - At least 1 uppercase letter
 * - At least 1 lowercase letter
 * - At least 1 number
 * - At least 1 special character
 * 
 * @param {string} password - Password to validate
 * @returns {Object} { valid: boolean, errors: string[], strength: number }
 */
export function validatePassword(password) {
  const errors = []
  let strength = 0
  
  if (!password) {
    return { valid: false, errors: ['Le mot de passe est requis'], strength: 0 }
  }
  
  // Longueur minimum
  if (password.length < 8) {
    errors.push('Au moins 8 caractères')
  } else {
    strength += 1
  }
  
  // Majuscules
  if (!/[A-Z]/.test(password)) {
    errors.push('Au moins 1 majuscule')
  } else {
    strength += 1
  }
  
  // Minuscules
  if (!/[a-z]/.test(password)) {
    errors.push('Au moins 1 minuscule')
  } else {
    strength += 1
  }
  
  // Chiffres
  if (!/[0-9]/.test(password)) {
    errors.push('Au moins 1 chiffre')
  } else {
    strength += 1
  }
  
  // Caractères spéciaux
  if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
    errors.push('Au moins 1 caractère spécial (!@#$%^&*...)')
  } else {
    strength += 1
  }
  
  // Bonus pour longueur > 12
  if (password.length >= 12) {
    strength += 1
  }
  
  // Vérifier mots de passe communs faibles
  const commonPasswords = [
    'password', 'password123', '12345678', 'qwerty', 'azerty',
    'admin', 'letmein', 'welcome', 'monkey', 'dragon',
    'master', 'abc123', '123456789', '1234567890', 'password1'
  ]
  
  if (commonPasswords.includes(password.toLowerCase())) {
    errors.push('Mot de passe trop commun, choisissez-en un plus unique')
    strength = 0
  }
  
  return {
    valid: errors.length === 0,
    errors,
    strength: Math.min(strength, 5), // Max 5
    strengthLabel: getStrengthLabel(strength)
  }
}

/**
 * Get password strength label
 * @param {number} strength - Strength score (0-5)
 * @returns {string} Label
 */
function getStrengthLabel(strength) {
  if (strength <= 1) return 'Très faible'
  if (strength === 2) return 'Faible'
  if (strength === 3) return 'Moyen'
  if (strength === 4) return 'Fort'
  return 'Très fort'
}

/**
 * Get password strength color
 * @param {number} strength - Strength score (0-5)
 * @returns {string} Color class
 */
export function getStrengthColor(strength) {
  if (strength <= 1) return '#ff3c3c' // Rouge
  if (strength === 2) return '#ff9500' // Orange
  if (strength === 3) return '#ffcc00' // Jaune
  if (strength === 4) return '#34c759' // Vert
  return '#00c853' // Vert foncé
}

/**
 * Validates username format (alphanumeric + underscore, 3-20 chars)
 * @param {string} username - Username to validate
 * @returns {Object} { valid: boolean, error: string }
 */
export function validateUsername(username) {
  if (!username) {
    return { valid: false, error: 'Le nom d\'utilisateur est requis' }
  }
  
  if (username.length < 3) {
    return { valid: false, error: 'Au moins 3 caractères' }
  }
  
  if (username.length > 20) {
    return { valid: false, error: 'Maximum 20 caractères' }
  }
  
  if (!/^[a-zA-Z0-9_]+$/.test(username)) {
    return { valid: false, error: 'Uniquement lettres, chiffres et underscore (_)' }
  }
  
  return { valid: true, error: '' }
}

/**
 * Validates phone number format (international or french)
 * @param {string} phone - Phone number to validate
 * @returns {boolean} True if valid or empty
 */
export function isValidPhone(phone) {
  if (!phone || phone.trim() === '') return true // Optional field
  
  // Remove spaces, dots, dashes
  const cleaned = phone.replace(/[\s.\-()]/g, '')
  
  // French format: 0612345678 or +33612345678
  const frenchRegex = /^(\+33|0)[1-9]\d{8}$/
  
  // International format: +[country code][number]
  const intlRegex = /^\+\d{1,3}\d{4,14}$/
  
  return frenchRegex.test(cleaned) || intlRegex.test(cleaned)
}
