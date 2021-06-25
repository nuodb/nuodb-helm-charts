# NuoDB Storage Manager (SM) Journal in Kubernetes

## Introduction

The storage manager, SM for short, is responsible for maintaining a complete copy of the database.
The atoms, database elements, are stored to either local disk, or in a separate volume.
These atoms are written by a specific module called the `archive`.
However, the archive doesn’t have any special mechanisms to satisfy the durability requirements.
In addition to the archive, NuoDB storage managers also maintain a write-ahead-log called the `journal`.
The journal has the onerous task of ensuring durability in the face of an unexpected process termination.
E.g., a power loss, machine meltdown, or an unexpected cloud outage.
The journal will make sure that the archive is reconstructed in a consistent state when coming back online.

For more info on NuoDB database journaling, please consult the official [NuoDB docs](https://doc.nuodb.com/nuodb/latest/database-administration/about-database-journaling/).

## Separating the Journal from the Archive

Since the `journal` has to write all commits durably to disk, the speed of the disk directly influences the commit latency.
To achieve the best performance, NuoDB recommends placing the `journal` on the fastest disk available.

By default, the `journal` is located in a subdirectory of the `archive`.
To achieve the best cost vs speed tradeoff, you can separate the journal from the archive.
To achieve that set the database helm template values `database.sm.noHotCopy.journalPath.enabled` and `database.sm.hotCopy.journalPath.enabled` to `true` and configure it with the desired persistence settings.

Kubernetes stateful sets volume mounts are immutable and as such, the setting can not be changed easily on an existing database.
More on upgrades below [Upgrading existing domains](#upgrading-existing-domains).

## Upgrading existing domains
TODO