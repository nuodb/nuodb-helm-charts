# Install NuoDB License on Kubernetes

When using the Helm package manager to deploy NuoDB, a license is not installed by default.

To view the details such as date of expiry (`expires`), name of the license holder (`holder`), and the type of license (`type`) , run:

```
  nuocmd get effective-license
```

For example, if the *Enterprise License* is installed:

```
nuocmd --show-json get effective-license

{
  "decodedLicense": {
    "expires": "2025-01-31T15:33:40.465134",
    "holder": "3DS Solutions",
    "type": "ENTERPRISE"
  },
  "effectiveForDomain": true,
  "encodedLicense":
    "-----BEGIN LICENSE-----
    <base64-encoded data>
    -----END LICENSE-----"
}
```
To install a license or to upgrade an existing license, choose one of the following options:

* Redeploy the Admin Processes (APs) using Helm [RECOMMENDED]
* Use the `nuocmd set license` command

## Redeploy the Admin Processes (APs) using Helm

### To install an *Enterprise License* or a *Limited License* for a new Admin domain

1. Contact <nuodb.support@3ds.com> to obtain a new Limited License key or an Enterprise License key.

2. Make a copy of the admin chart’s `values.yaml` file.

   ```
   helm show values nuodb/admin > admin-values.yaml
   ```

3. Edit the `admin-values.yaml` file.

   Replace the contents of `configFiles:` with the contents of the license file. For example:

   ```
   configFiles:
    nuodb.lic: |-
    #   "PUT YOUR BASE64 ENCODED LICENSE CONTENT HERE"
    -----BEGIN LICENSE-----
    <base64-encoded data>
    -----END LICENSE-----
   ```

   > **NOTE**:
   > Paste the entire license file, including the lines `-----BEGIN LICENSE-----` and `-----END LICENSE-----`.

4. Save the changes to the `admin-values.yaml` file.

5. Re-install the admin chart specifying `--values admin-values.yaml`.

   ```
   helm upgrade --install <RELEASE_NAME> nuodb/admin --values admin-values.yaml
   ```

6. Check the details of the updated license.

   ```
   nuocmd --show-json get effective-license
   ```

### To install an *Enterprise License* for an existing Admin domain

1. Contact <nuodb.support@3ds.com> to get a new *Enterprise License*.

2. Make a copy of the admin chart’s Helm values.

   ```
   helm get values --all --output=yaml <RELEASE_NAME> admin-values.yaml
   ```
3. Edit the `admin-values.yaml` file.

   Replace the contents of `configFiles:` with the following:

   ```
   configFiles:
    nuodb.lic: |-
    #   "PUT YOUR BASE64 ENCODED LICENSE CONTENT HERE"
    -----BEGIN LICENSE-----
    <base64-encoded data>
    -----END LICENSE-----
   ```
   > **NOTE**:
   > Paste the entire license file, including the lines `-----BEGIN LICENSE-----` and `-----END LICENSE-----`.

4. Save the changes to the `admin-values.yaml` file.

5. Re-run the admin chart specifying `--values admin-values.yaml`.

   ```
   helm upgrade <RELEASE_NAME> nuodb/admin --values admin-values.yaml
   ```
   This step will restart the AP pods one at a time without affecting any running applications or databases.

6. Check the details of the updated license.

   ```
   nuocmd --show-json get effective-license
   ```

   > **TIP**:
   >Use a version control software to track the changes to the `admin-values.yaml` file.

## Use the nuocmd set license command
To install a NuoDB *Enterprise License* for an existing Admin domain using nuocmd, invoke `nuocmd` on an AP running in Kubernetes.

   ```
   kubectl cp <nuodb.lic path on local host> <AP-pod-name>:/tmp/nuodb.lic
   kubectl exec <AP-pod-name> -- nuocmd set license --license-file /tmp/nuodb.lic
   ```
> **NOTE**:
> Since the license is stored in the key-value store of the Raft state and is replicated automatically to all APs, run `nuocmd set license` on any one AP in the domain.
