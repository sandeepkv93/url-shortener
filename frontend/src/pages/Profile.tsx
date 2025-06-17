const Profile = () => {
  return (
    <div className="px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Profile</h1>
        <p className="mt-2 text-gray-600">Manage your account settings and preferences.</p>
      </div>

      <div className="max-w-2xl">
        <div className="card mb-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Account Information</h3>
          <div className="space-y-4">
            <div>
              <label htmlFor="name" className="label">
                Full Name
              </label>
              <input
                type="text"
                id="name"
                name="name"
                className="input"
                placeholder="John Doe"
              />
            </div>
            <div>
              <label htmlFor="email" className="label">
                Email Address
              </label>
              <input
                type="email"
                id="email"
                name="email"
                className="input"
                placeholder="john@example.com"
              />
            </div>
            <div className="flex justify-end">
              <button className="btn-primary">
                Save Changes
              </button>
            </div>
          </div>
        </div>

        <div className="card mb-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Change Password</h3>
          <div className="space-y-4">
            <div>
              <label htmlFor="current-password" className="label">
                Current Password
              </label>
              <input
                type="password"
                id="current-password"
                name="current-password"
                className="input"
              />
            </div>
            <div>
              <label htmlFor="new-password" className="label">
                New Password
              </label>
              <input
                type="password"
                id="new-password"
                name="new-password"
                className="input"
              />
            </div>
            <div>
              <label htmlFor="confirm-password" className="label">
                Confirm New Password
              </label>
              <input
                type="password"
                id="confirm-password"
                name="confirm-password"
                className="input"
              />
            </div>
            <div className="flex justify-end">
              <button className="btn-primary">
                Update Password
              </button>
            </div>
          </div>
        </div>

        <div className="card">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Account Settings</h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">Email Notifications</h4>
                <p className="text-sm text-gray-500">Receive notifications about your URLs</p>
              </div>
              <button
                type="button"
                className="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent bg-gray-200 transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                role="switch"
                aria-checked="false"
              >
                <span className="translate-x-0 inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"></span>
              </button>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-gray-900">Public Profile</h4>
                <p className="text-sm text-gray-500">Make your profile visible to others</p>
              </div>
              <button
                type="button"
                className="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent bg-gray-200 transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                role="switch"
                aria-checked="false"
              >
                <span className="translate-x-0 inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out"></span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Profile