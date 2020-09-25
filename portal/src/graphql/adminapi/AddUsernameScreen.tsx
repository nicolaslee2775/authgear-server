import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { useNavigate, useParams } from "react-router-dom";
import { TextField } from "@fluentui/react";
import { Context, FormattedMessage } from "@oursky/react-messageformat";
import deepEqual from "deep-equal";

import UserDetailCommandBar from "./UserDetailCommandBar";
import NavBreadcrumb from "../../NavBreadcrumb";
import NavigationBlockerDialog from "../../NavigationBlockerDialog";
import ButtonWithLoading from "../../ButtonWithLoading";
import ShowError from "../../ShowError";
import { useCreateLoginIDIdentityMutation } from "./mutations/createIdentityMutation";
import { useTextField } from "../../hook/useInput";
import {
  defaultFormatErrorMessageList,
  Violation,
} from "../../util/validation";
import { parseError } from "../../util/error";

import styles from "./AddUsernameScreen.module.scss";

const AddUsernameScreen: React.FC = function AddUsernameScreen() {
  const { userID } = useParams();
  const navigate = useNavigate();

  const {
    createIdentity,
    loading: creatingIdentity,
    error: createIdentityError,
  } = useCreateLoginIDIdentityMutation(userID);
  const { renderToString } = useContext(Context);

  const [submittedForm, setSubmittedForm] = useState<boolean>(false);

  const { value: username, onChange: onUsernameChange } = useTextField("");

  const navBreadcrumbItems = useMemo(() => {
    return [
      { to: "../../..", label: <FormattedMessage id="UsersScreen.title" /> },
      { to: "../", label: <FormattedMessage id="UserDetailsScreen.title" /> },
      { to: ".", label: <FormattedMessage id="AddUsernameScreen.title" /> },
    ];
  }, []);

  const screenState = useMemo(
    () => ({
      username,
    }),
    [username]
  );

  const isFormModified = useMemo(() => {
    return !deepEqual({ username: "" }, screenState);
  }, [screenState]);

  const onAddClicked = useCallback(() => {
    createIdentity({ key: "username", value: username })
      .then((identity) => {
        if (identity != null) {
          setSubmittedForm(true);
        }
      })
      .catch(() => {});
  }, [username, createIdentity]);

  useEffect(() => {
    if (submittedForm) {
      navigate("../#connected-identities");
    }
  }, [submittedForm, navigate]);

  const { errorMessage, unhandledViolations } = useMemo(() => {
    const violations = parseError(createIdentityError);
    const usernameFieldErrorMessages: string[] = [];
    const unhandledViolations: Violation[] = [];
    for (const violation of violations) {
      if (violation.kind === "Invalid" || violation.kind === "format") {
        usernameFieldErrorMessages.push(
          renderToString("AddUsernameScreen.error.invalid-username")
        );
      } else if (violation.kind === "DuplicatedIdentity") {
        usernameFieldErrorMessages.push(
          renderToString("AddUsernameScreen.error.duplicated-username")
        );
      } else {
        unhandledViolations.push(violation);
      }
    }

    const errorMessage = {
      username: defaultFormatErrorMessageList(usernameFieldErrorMessages),
    };

    return { errorMessage, unhandledViolations };
  }, [createIdentityError, renderToString]);

  return (
    <div className={styles.root}>
      <UserDetailCommandBar />
      <NavBreadcrumb className={styles.breadcrumb} items={navBreadcrumbItems} />
      <section className={styles.content}>
        {unhandledViolations.length > 0 && (
          <ShowError error={createIdentityError} />
        )}
        <NavigationBlockerDialog
          blockNavigation={!submittedForm && isFormModified}
        />
        <TextField
          className={styles.usernameField}
          label={renderToString("AddUsernameScreen.username.label")}
          value={username}
          onChange={onUsernameChange}
          errorMessage={errorMessage.username}
        />
        <ButtonWithLoading
          onClick={onAddClicked}
          disabled={!isFormModified}
          labelId="add"
          loading={creatingIdentity}
        />
      </section>
    </div>
  );
};

export default AddUsernameScreen;