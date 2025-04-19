#!/bin/bash

YQ_IMAGE="mikefarah/yq:latest"

template_source=""
values_args=()

# --- Parse command line arguments ---
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --template)
      template_source="$2"
      shift # past argument
      shift # past value
      ;;
    --values)
      # Collect all arguments after --values until next -- argument or end
      shift # past --values
      while [[ $# -gt 0 ]] && [[ ! "$1" =~ ^-- ]]; do
        values_args+=("$1")
        shift # past value argument
      done
      ;;
    *) # Unknown option
      echo "Error: Unknown option $1" >&2
      exit 1
      ;;
  esac
done

# --- Get template content ---
current_yaml=""
if [[ -z "$template_source" ]]; then
  # If no template provided, start with empty YAML
  current_yaml="{}"
elif [[ "$template_source" =~ ^https?:// ]]; then
  # Download from URL
  # Use curl, exit if fails
  if ! command -v curl &> /dev/null; then
      echo "Error: curl command required to download URL template." >&2
      exit 1
  fi
  template_content=$(curl -sfL "$template_source")
  if [[ $? -ne 0 ]]; then
    echo "Error: Failed to download template from URL: $template_source" >&2
    exit 1
  fi
  # Check if downloaded content is empty
  if [[ -z "$template_content" ]]; then
       current_yaml="{}"
  else
       current_yaml="$template_content"
  fi

elif [[ -f "$template_source" ]]; then
  # Read from local file
  current_yaml=$(cat "$template_source")
  # Check if file content is empty
  if [[ -z "$current_yaml" ]]; then
       current_yaml="{}"
  fi
else
  # Invalid template source
  echo "Error: Invalid template source '$template_source'. Please provide valid file path or HTTP/HTTPS URL." >&2
  exit 1
fi


# --- Check if Docker is available ---
if ! command -v docker &> /dev/null; then
    echo "Error: docker command required to run yq." >&2
    exit 1
fi
# Try pulling or verifying yq image (helps catch issues early)
docker pull $YQ_IMAGE > /dev/null

# --- Apply values ---
if [[ ${#values_args[@]} -gt 0 ]]; then
  for val_arg in "${values_args[@]}"; do
    # Parse key=value
    if [[ ! "$val_arg" =~ ^([^=]+)=(.*)$ ]]; then
      continue
    fi

    # BASH_REMATCH is result array from =~ operator
    yaml_path="${BASH_REMATCH[1]}"
    raw_value="${BASH_REMATCH[2]}"

    # Prepare yq value (try handling basic types, otherwise treat as string)
    yq_value=""
    if [[ "$raw_value" == "true" || "$raw_value" == "false" || "$raw_value" == "null" ]]; then
      yq_value="$raw_value"
    # Check if integer or float (simple regex)
    elif [[ "$raw_value" =~ ^-?[0-9]+(\.[0-9]+)?$ ]]; then
       # If value starts with 0 but isn't 0 itself and has no decimal point, force string to prevent octal interpretation
       if [[ "$raw_value" =~ ^0[0-9]+$ ]]; then
           # Need to escape internal double quotes
           escaped_value=$(echo "$raw_value" | sed 's/"/\\"/g')
           yq_value="\"$escaped_value\""
       else
           yq_value="$raw_value"
       fi
    else
      # Treat as string, need to escape internal double quotes
      escaped_value=$(echo "$raw_value" | sed 's/"/\\"/g')
      yq_value="\"$escaped_value\""
    fi

    # Build yq expression
    yq_expression=".$yaml_path = $yq_value"

    # Apply update via docker run yq
    # Pass current YAML via stdin to yq, get stdout as new YAML
    # Use <<< for here-string input to avoid temp files
    new_yaml=$(docker run --rm -i "$YQ_IMAGE" "$yq_expression" <<< "$current_yaml")
    yq_exit_code=$?

    if [[ $yq_exit_code -ne 0 ]]; then
      echo "Error: yq execution failed (exit code: $yq_exit_code). Expression: '$yq_expression'" >&2
      # Could output yq error message, but requires more complex docker run call to capture stderr
      exit 1
    fi
    current_yaml="$new_yaml"
  done
fi

# --- Output final result ---
printf "%s\n" "$current_yaml"

exit 0