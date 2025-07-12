import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import { UpdateProfileRequest, ChangePasswordRequest } from '@/types/auth'
import { PageLoading, ButtonLoading } from '@/components/common/Loading'
import {
  User,
  Mail,
  Lock,
  Eye,
  EyeOff,
  Save,
  AlertCircle,
  CheckCircle,
  Settings,
  Bell,
  Globe,
  ArrowLeft,
  Trash2,
  Download
} from 'lucide-react'

// Validation schemas
const profileSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Please enter a valid email address')
})

const passwordSchema = z.object({
  currentPassword: z.string().min(1, 'Current password is required'),
  newPassword: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string().min(1, 'Please confirm your password')
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"]
})

type ProfileFormData = z.infer<typeof profileSchema>
type PasswordFormData = z.infer<typeof passwordSchema>

interface UserSettings {
  emailNotifications: boolean
  publicProfile: boolean
  marketingEmails: boolean
  twoFactorAuth: boolean
}

const Profile = () => {
  const { user, updateProfile, changePassword, isLoading: authLoading, logout } = useAuth()
  const [isUpdatingProfile, setIsUpdatingProfile] = useState(false)
  const [isChangingPassword, setIsChangingPassword] = useState(false)
  const [showCurrentPassword, setShowCurrentPassword] = useState(false)
  const [showNewPassword, setShowNewPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [profileSuccess, setProfileSuccess] = useState('')
  const [passwordSuccess, setPasswordSuccess] = useState('')
  const [settings, setSettings] = useState<UserSettings>({
    emailNotifications: true,
    publicProfile: false,
    marketingEmails: false,
    twoFactorAuth: false
  })

  const {
    register: registerProfile,
    handleSubmit: handleProfileSubmit,
    formState: { errors: profileErrors },
    setValue: setProfileValue,
    reset: resetProfile
  } = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: user?.name || '',
      email: user?.email || ''
    }
  })

  const {
    register: registerPassword,
    handleSubmit: handlePasswordSubmit,
    formState: { errors: passwordErrors },
    reset: resetPassword
  } = useForm<PasswordFormData>({
    resolver: zodResolver(passwordSchema)
  })

  useEffect(() => {
    if (user) {
      setProfileValue('name', user.name)
      setProfileValue('email', user.email)
    }
  }, [user, setProfileValue])

  const handleUpdateProfile = async (data: ProfileFormData) => {
    if (!user) return

    setIsUpdatingProfile(true)
    setProfileSuccess('')

    try {
      const updateData: UpdateProfileRequest = {
        name: data.name,
        email: data.email
      }
      
      await updateProfile(updateData)
      setProfileSuccess('Profile updated successfully!')
      
      // Clear success message after 5 seconds
      setTimeout(() => setProfileSuccess(''), 5000)
    } catch (error: any) {
      console.error('Failed to update profile:', error)
    } finally {
      setIsUpdatingProfile(false)
    }
  }

  const handleChangePassword = async (data: PasswordFormData) => {
    setIsChangingPassword(true)
    setPasswordSuccess('')

    try {
      const changeData: ChangePasswordRequest = {
        currentPassword: data.currentPassword,
        newPassword: data.newPassword,
        confirmPassword: data.confirmPassword
      }
      
      await changePassword(changeData)
      setPasswordSuccess('Password changed successfully!')
      resetPassword()
      
      // Clear success message after 5 seconds
      setTimeout(() => setPasswordSuccess(''), 5000)
    } catch (error: any) {
      console.error('Failed to change password:', error)
    } finally {
      setIsChangingPassword(false)
    }
  }

  const handleSettingToggle = (setting: keyof UserSettings) => {
    setSettings(prev => ({
      ...prev,
      [setting]: !prev[setting]
    }))
    // In a real app, this would save to the backend
    console.log(`Toggled ${setting}:`, !settings[setting])
  }

  const handleDeleteAccount = () => {
    if (window.confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
      // In a real app, this would delete the account
      console.log('Delete account requested')
    }
  }

  const handleExportData = () => {
    // In a real app, this would export user data
    console.log('Export data requested')
  }

  if (!user) {
    return (
      <div className="px-4 py-8">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Please log in to view your profile</h1>
          <Link to="/login" className="btn-primary">
            Log In
          </Link>
        </div>
      </div>
    )
  }

  if (authLoading) {
    return <PageLoading message="Loading profile..." />
  }

  const ToggleSwitch = ({ 
    enabled, 
    onChange, 
    disabled = false 
  }: { 
    enabled: boolean
    onChange: () => void
    disabled?: boolean 
  }) => (
    <button
      type="button"
      onClick={onChange}
      disabled={disabled}
      className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 ${
        enabled ? 'bg-primary-600' : 'bg-gray-200'
      } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
      role="switch"
      aria-checked={enabled}
    >
      <span
        className={`${
          enabled ? 'translate-x-5' : 'translate-x-0'
        } inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out`}
      />
    </button>
  )

  return (
    <div className="px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center space-x-3 mb-4">
          <Link 
            to="/dashboard"
            className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
          >
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Profile Settings</h1>
            <p className="mt-2 text-gray-600">Manage your account settings and preferences.</p>
          </div>
        </div>
      </div>

      <div className="max-w-2xl space-y-6">
        {/* Account Information */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center space-x-3 mb-6">
            <div className="h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
              <User className="h-5 w-5 text-blue-600" />
            </div>
            <h3 className="text-lg font-medium text-gray-900">Account Information</h3>
          </div>

          {profileSuccess && (
            <div className="mb-4 p-3 bg-green-50 border border-green-200 rounded-md">
              <div className="flex items-center">
                <CheckCircle className="h-5 w-5 text-green-500 mr-2" />
                <p className="text-sm text-green-700">{profileSuccess}</p>
              </div>
            </div>
          )}

          <form onSubmit={handleProfileSubmit(handleUpdateProfile)} className="space-y-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                Full Name
              </label>
              <input
                type="text"
                id="name"
                className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                  profileErrors.name 
                    ? 'border-red-300 bg-red-50' 
                    : 'border-gray-300 hover:border-gray-400'
                }`}
                {...registerProfile('name')}
              />
              {profileErrors.name && (
                <p className="mt-1 text-sm text-red-600 flex items-center">
                  <AlertCircle className="h-4 w-4 mr-1" />
                  {profileErrors.name.message}
                </p>
              )}
            </div>

            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                Email Address
              </label>
              <input
                type="email"
                id="email"
                className={`w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                  profileErrors.email 
                    ? 'border-red-300 bg-red-50' 
                    : 'border-gray-300 hover:border-gray-400'
                }`}
                {...registerProfile('email')}
              />
              {profileErrors.email && (
                <p className="mt-1 text-sm text-red-600 flex items-center">
                  <AlertCircle className="h-4 w-4 mr-1" />
                  {profileErrors.email.message}
                </p>
              )}
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                disabled={isUpdatingProfile}
                className={`inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white transition-colors ${
                  isUpdatingProfile
                    ? 'bg-gray-400 cursor-not-allowed'
                    : 'bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500'
                }`}
              >
                {isUpdatingProfile ? (
                  <>
                    <ButtonLoading className="mr-2" />
                    Saving...
                  </>
                ) : (
                  <>
                    <Save className="h-4 w-4 mr-2" />
                    Save Changes
                  </>
                )}
              </button>
            </div>
          </form>
        </div>

        {/* Change Password */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center space-x-3 mb-6">
            <div className="h-8 w-8 bg-yellow-100 rounded-lg flex items-center justify-center">
              <Lock className="h-5 w-5 text-yellow-600" />
            </div>
            <h3 className="text-lg font-medium text-gray-900">Change Password</h3>
          </div>

          {passwordSuccess && (
            <div className="mb-4 p-3 bg-green-50 border border-green-200 rounded-md">
              <div className="flex items-center">
                <CheckCircle className="h-5 w-5 text-green-500 mr-2" />
                <p className="text-sm text-green-700">{passwordSuccess}</p>
              </div>
            </div>
          )}

          <form onSubmit={handlePasswordSubmit(handleChangePassword)} className="space-y-4">
            <div>
              <label htmlFor="currentPassword" className="block text-sm font-medium text-gray-700 mb-1">
                Current Password
              </label>
              <div className="relative">
                <input
                  type={showCurrentPassword ? 'text' : 'password'}
                  id="currentPassword"
                  className={`w-full px-3 py-2 pr-10 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                    passwordErrors.currentPassword 
                      ? 'border-red-300 bg-red-50' 
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                  {...registerPassword('currentPassword')}
                />
                <button
                  type="button"
                  onClick={() => setShowCurrentPassword(!showCurrentPassword)}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
                >
                  {showCurrentPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {passwordErrors.currentPassword && (
                <p className="mt-1 text-sm text-red-600 flex items-center">
                  <AlertCircle className="h-4 w-4 mr-1" />
                  {passwordErrors.currentPassword.message}
                </p>
              )}
            </div>

            <div>
              <label htmlFor="newPassword" className="block text-sm font-medium text-gray-700 mb-1">
                New Password
              </label>
              <div className="relative">
                <input
                  type={showNewPassword ? 'text' : 'password'}
                  id="newPassword"
                  className={`w-full px-3 py-2 pr-10 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                    passwordErrors.newPassword 
                      ? 'border-red-300 bg-red-50' 
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                  {...registerPassword('newPassword')}
                />
                <button
                  type="button"
                  onClick={() => setShowNewPassword(!showNewPassword)}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
                >
                  {showNewPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {passwordErrors.newPassword && (
                <p className="mt-1 text-sm text-red-600 flex items-center">
                  <AlertCircle className="h-4 w-4 mr-1" />
                  {passwordErrors.newPassword.message}
                </p>
              )}
            </div>

            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-1">
                Confirm New Password
              </label>
              <div className="relative">
                <input
                  type={showConfirmPassword ? 'text' : 'password'}
                  id="confirmPassword"
                  className={`w-full px-3 py-2 pr-10 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors ${
                    passwordErrors.confirmPassword 
                      ? 'border-red-300 bg-red-50' 
                      : 'border-gray-300 hover:border-gray-400'
                  }`}
                  {...registerPassword('confirmPassword')}
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
                >
                  {showConfirmPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {passwordErrors.confirmPassword && (
                <p className="mt-1 text-sm text-red-600 flex items-center">
                  <AlertCircle className="h-4 w-4 mr-1" />
                  {passwordErrors.confirmPassword.message}
                </p>
              )}
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                disabled={isChangingPassword}
                className={`inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white transition-colors ${
                  isChangingPassword
                    ? 'bg-gray-400 cursor-not-allowed'
                    : 'bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500'
                }`}
              >
                {isChangingPassword ? (
                  <>
                    <ButtonLoading className="mr-2" />
                    Updating...
                  </>
                ) : (
                  <>
                    <Lock className="h-4 w-4 mr-2" />
                    Update Password
                  </>
                )}
              </button>
            </div>
          </form>
        </div>

        {/* Account Settings */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center space-x-3 mb-6">
            <div className="h-8 w-8 bg-green-100 rounded-lg flex items-center justify-center">
              <Settings className="h-5 w-5 text-green-600" />
            </div>
            <h3 className="text-lg font-medium text-gray-900">Preferences</h3>
          </div>

          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Bell className="h-5 w-5 text-gray-400" />
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Email Notifications</h4>
                  <p className="text-sm text-gray-500">Receive notifications about your URLs and account activity</p>
                </div>
              </div>
              <ToggleSwitch
                enabled={settings.emailNotifications}
                onChange={() => handleSettingToggle('emailNotifications')}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Globe className="h-5 w-5 text-gray-400" />
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Public Profile</h4>
                  <p className="text-sm text-gray-500">Make your profile visible to other users</p>
                </div>
              </div>
              <ToggleSwitch
                enabled={settings.publicProfile}
                onChange={() => handleSettingToggle('publicProfile')}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Mail className="h-5 w-5 text-gray-400" />
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Marketing Emails</h4>
                  <p className="text-sm text-gray-500">Receive updates about new features and promotions</p>
                </div>
              </div>
              <ToggleSwitch
                enabled={settings.marketingEmails}
                onChange={() => handleSettingToggle('marketingEmails')}
              />
            </div>

            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Lock className="h-5 w-5 text-gray-400" />
                <div>
                  <h4 className="text-sm font-medium text-gray-900">Two-Factor Authentication</h4>
                  <p className="text-sm text-gray-500">Add an extra layer of security to your account</p>
                </div>
              </div>
              <ToggleSwitch
                enabled={settings.twoFactorAuth}
                onChange={() => handleSettingToggle('twoFactorAuth')}
                disabled={true}
              />
            </div>
          </div>
        </div>

        {/* Data & Privacy */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Data & Privacy</h3>
          <div className="space-y-4">
            <button
              onClick={handleExportData}
              className="flex items-center justify-between w-full p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div className="flex items-center space-x-3">
                <Download className="h-5 w-5 text-gray-600" />
                <div className="text-left">
                  <div className="font-medium text-gray-900">Export Your Data</div>
                  <div className="text-sm text-gray-500">Download all your URLs and analytics data</div>
                </div>
              </div>
            </button>

            <button
              onClick={handleDeleteAccount}
              className="flex items-center justify-between w-full p-4 border border-red-200 rounded-lg hover:bg-red-50 transition-colors text-red-600"
            >
              <div className="flex items-center space-x-3">
                <Trash2 className="h-5 w-5" />
                <div className="text-left">
                  <div className="font-medium">Delete Account</div>
                  <div className="text-sm text-red-500">Permanently delete your account and all data</div>
                </div>
              </div>
            </button>
          </div>
        </div>

        {/* Account Info */}
        <div className="bg-gray-50 rounded-lg p-4">
          <h4 className="text-sm font-medium text-gray-900 mb-2">Account Details</h4>
          <div className="text-sm text-gray-600 space-y-1">
            <p><strong>User ID:</strong> {user.id}</p>
            <p><strong>Member since:</strong> {new Date(user.createdAt).toLocaleDateString()}</p>
            <p><strong>Last updated:</strong> {new Date(user.updatedAt).toLocaleDateString()}</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Profile